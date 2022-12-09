package main

import (
	"github.com/davecheney/m/m"
	"github.com/davecheney/m/mastodon"
	"gorm.io/gorm"
)

type AutoMigrateCmd struct {
}

func (a *AutoMigrateCmd) Run(ctx *Context) error {
	db, err := gorm.Open(ctx.Dialector, &ctx.Config)
	if err != nil {
		return err
	}

	return db.AutoMigrate(
		&m.Account{}, &m.AccountList{}, &m.LocalAccount{},
		&m.Activity{},
		&mastodon.Application{},
		&m.Conversation{},
		&mastodon.ClientFilter{},
		&m.Instance{}, &m.InstanceRule{},
		&m.Marker{},
		&m.Notification{},
		&m.Status{},
		&m.Token{},
	)
}
