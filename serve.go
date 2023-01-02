package main

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/davecheney/pub/activitypub"
	"github.com/davecheney/pub/internal/group"
	"github.com/davecheney/pub/internal/models"
	"github.com/davecheney/pub/mastodon"
	"github.com/davecheney/pub/oauth"
	"github.com/davecheney/pub/wellknown"
	"gorm.io/gorm"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ServeCmd struct {
	Addr             string `help:"address to listen" default:"127.0.0.1:9999"`
	DebugPrintRoutes bool   `help:"print routes to stdout on startup"`
	LogHTTP          bool   `help:"log HTTP requests"`
}

func (s *ServeCmd) Run(ctx *Context) error {
	db, err := gorm.Open(ctx.Dialector, &ctx.Config)
	if err != nil {
		return err
	}

	if err := configureDB(db); err != nil {
		return err
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	if s.LogHTTP {
		r.Use(middleware.Logger)
	}
	r.Use(setDBMiddleware(db))

	r.Route("/api", func(r chi.Router) {
		m := mastodon.NewService(db)
		instance := m.Instances()
		r.Route("/v1", func(r chi.Router) {
			r.Post("/apps", m.Applications().Create)
			r.Route("/accounts", func(r chi.Router) {
				accounts := m.Accounts()
				r.Get("/verify_credentials", accounts.VerifyCredentials)
				r.Patch("/update_credentials", accounts.Update)
				r.Get("/relationships", m.Relationships().Show)
				r.Get("/filters", m.Filters().Index)
				r.Get("/{id}", accounts.Show)
				r.Get("/{id}/lists", m.Lists().ShowListMembership)
				r.Get("/{id}/statuses", accounts.StatusesShow)
				r.Post("/{id}/follow", m.Relationships().Create)
				r.Get("/{id}/followers", accounts.FollowersShow)
				r.Get("/{id}/following", accounts.FollowingShow)
				r.Post("/{id}/unfollow", m.Relationships().Destroy)
				r.Post("/{id}/mute", m.Mutes().Create)
				r.Post("/{id}/unmute", m.Mutes().Destroy)
				r.Post("/{id}/block", m.Blocks().Create)
				r.Post("/{id}/unblock", m.Blocks().Destroy)
			})
			r.Get("/blocks", m.Blocks().Index)
			r.Get("/conversations", m.Conversations().Index)
			r.Get("/custom_emojis", m.Emojis().Index)
			r.Get("/directory", m.Directory().Index)
			r.Get("/filters", m.Filters().Index)
			r.Get("/lists", m.Lists().Index)
			r.Post("/lists", m.Lists().Create)
			r.Get("/lists/{id}", m.Lists().Show)
			r.Get("/lists/{id}/accounts", m.Lists().ViewMembers)
			r.Post("/lists/{id}/accounts", m.Lists().AddMembers)
			r.Delete("/lists/{id}/accounts", m.Lists().RemoveMembers)
			r.Get("/instance", instance.IndexV1)
			r.Get("/instance/", instance.IndexV1) // sigh
			r.Get("/instance/peers", mastodon.InstancesPeersShow)
			r.Get("/instance/activity", instance.ActivityShow)
			r.Get("/instance/domain_blocks", instance.DomainBlocksShow)
			r.Get("/markers", m.Markers().Index)
			r.Post("/markers", m.Markers().Create)
			r.Get("/mutes", m.Mutes().Index)
			r.Get("/notifications", m.Notifications().Index)

			r.Post("/statuses", m.Statuses().Create)
			r.Get("/statuses/{id}/context", m.Contexts().Show)
			r.Post("/statuses/{id}/favourite", m.Favourites().Create)
			r.Post("/statuses/{id}/unfavourite", m.Favourites().Destroy)
			r.Get("/statuses/{id}/favourited_by", m.Favourites().Show)
			r.Get("/statuses/{id}", m.Statuses().Show)
			r.Delete("/statuses/{id}", m.Statuses().Destroy)
			r.Route("/timelines", func(r chi.Router) {
				timelines := m.Timelines()
				r.Get("/home", timelines.Home)
				r.Get("/public", mastodon.TimelinesPublic)
			})

		})
		r.Route("/v2", func(r chi.Router) {
			r.Get("/instance", instance.IndexV2)
			r.Get("/search", mastodon.SearchIndex)
		})
		r.Route("/nodeinfo", func(r chi.Router) {
			r.Get("/2.0", wellknown.NodeInfoShow)
		})
	})

	ap := activitypub.NewService(db)
	getKey := func(keyID string) (crypto.PublicKey, error) {
		actorId := trimKeyId(keyID)
		var instance models.Instance
		if err := db.Joins("Admin").Preload("Admin.Actor").First(&instance, "admin_id is not null").Error; err != nil {
			return nil, err
		}
		fetcher := activitypub.NewRemoteActorFetcher(instance.Admin, db)
		actor, err := models.NewActors(db).FindOrCreate(actorId, fetcher.Fetch)
		if err != nil {
			return nil, err
		}
		return pemToPublicKey(actor.PublicKey)
	}
	r.Post("/inbox", ap.Inboxes(getKey).Create)

	r.Route("/oauth", func(r chi.Router) {
		r.Get("/authorize", oauth.AuthorizeNew)
		r.Post("/authorize", oauth.AuthorizeCreate)
		r.Post("/token", oauth.TokenCreate)
		r.Post("/revoke", oauth.TokenDestroy)
	})

	r.Route("/u/{username}", func(r chi.Router) {
		r.Get("/", activitypub.UsersShow)
		r.Post("/inbox", ap.Inboxes(getKey).Create)
		r.Get("/outbox", activitypub.OutboxIndex)
		r.Get("/followers", activitypub.FollowersIndex)
		r.Get("/following", activitypub.FollowingIndex)
		r.Get("/collections/{collection}", ap.Collections().Show)
	})

	r.Route("/.well-known", func(r chi.Router) {
		r.Get("/webfinger", wellknown.WebfingerShow)
		r.Get("/host-meta", wellknown.HostMetaIndex)
		r.Get("/nodeinfo", wellknown.NodeInfoIndex)
	})

	if s.DebugPrintRoutes {
		walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			route = strings.Replace(route, "/*/", "/", -1)
			fmt.Printf("%s %s\n", method, route)
			return nil
		}

		if err := chi.Walk(r, walkFunc); err != nil {
			fmt.Printf("Logging err: %s\n", err.Error())
		}
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	g := group.New(signalCtx)
	g.AddContext(func(ctx context.Context) error {
		fmt.Println("http.ListenAndServe", s.Addr, "started")
		defer fmt.Println("http.ListenAndServe", s.Addr, "stopped")
		svr := &http.Server{
			Addr:         s.Addr,
			Handler:      r,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		go func() {
			<-ctx.Done()
			svr.Shutdown(ctx)
		}()
		return svr.ListenAndServe()
	})
	relrp := activitypub.NewRelationshipRequestProcessor(db)
	g.Add(relrp.Run)
	reacrp := activitypub.NewReactionRequestProcessor(db)
	g.Add(reacrp.Run)

	return g.Wait()
}

func pemToPublicKey(key []byte) (crypto.PublicKey, error) {
	block, _ := pem.Decode(key)
	if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("pemToPublicKey: invalid pem type: %s", block.Type)
	}
	var publicKey interface{}
	var err error
	if publicKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return nil, fmt.Errorf("pemToPublicKey: parsepkixpublickey: %w", err)
	}
	return publicKey, nil
}

// trimKeyId removes the #main-key suffix from the key id.
func trimKeyId(id string) string {
	if i := strings.Index(id, "#"); i != -1 {
		return id[:i]
	}
	return id
}

func configureDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

func setDBMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "DB", db.WithContext(r.Context()))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
