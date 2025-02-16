package main

import (
	"github.com/davecheney/pub/internal/models"
	"gorm.io/gorm"
)

type AutoMigrateCmd struct {
	DisableForeignKeyConstraints bool `help:"disable foreign key constraints when creating tables."`
}

func (a *AutoMigrateCmd) Run(ctx *Context) error {
	ctx.Config.DisableForeignKeyConstraintWhenMigrating = a.DisableForeignKeyConstraints
	db, err := gorm.Open(ctx.Dialector, &ctx.Config)
	if err != nil {
		return err
	}

	return db.AutoMigrate(
		&models.ActivitypubRefresh{},
		&models.Actor{}, &models.ActorAttribute{},
		&models.Account{}, &models.AccountList{}, &models.AccountListMember{}, &models.AccountRole{}, &models.AccountMarker{},
		&models.Application{},
		&models.Conversation{},
		&models.Instance{}, &models.InstanceRule{},
		&models.Reaction{}, &models.ReactionRequest{},
		&models.Relationship{}, &models.RelationshipRequest{},
		// &models.Notification{},
		&models.Status{}, &models.StatusPoll{}, &models.StatusPollOption{}, &models.StatusAttachment{}, &models.StatusMention{}, &models.StatusTag{},
		&models.Tag{},
		&models.Token{},
	)
}
