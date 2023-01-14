package main

import (
	"strings"

	"github.com/alecthomas/kong"
	"gorm.io/driver/mysql"

	//"gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/go-sql-driver/mysql"
)

type Context struct {
	Debug bool

	gorm.Config
}

var cli struct {
	LogSQL bool   `help:"Log SQL queries."`
	DSN    string `help:"data source name" default:"pub:pub@tcp(localhost:3306)/pub"`
	Driver string `enum:"mysql,sqlite" help:"The database driver to use." required:"" default:"mysql"`

	AutoMigrate          AutoMigrateCmd          `cmd:"" help:"Automigrate the database."`
	CreateAccount        CreateAccountCmd        `cmd:"" help:"Create a new account."`
	CreateInstance       CreateInstanceCmd       `cmd:"" help:"Create a new instance."`
	DeleteAccount        DeleteAccountCmd        `cmd:"" help:"Delete an account."`
	Serve                ServeCmd                `cmd:"" help:"Serve a local web server."`
	SynchroniseFollowers SynchroniseFollowersCmd `cmd:"" help:"Synchronise followers."`
	Follow               FollowCmd               `cmd:"" help:"Follow an object."`
}

func main() {
	ctx := kong.Parse(&cli)

	var dialect gorm.Dialector

	switch cli.Driver {
	case "mysql":
		dialect = mysql.New(mysql.Config{
			DSN:                       mergeOptions(cli.DSN, "charset=utf8mb4&parseTime=True&loc=Local"),
			SkipInitializeWithVersion: false, // auto configure based on currently MySQL version

		})
	case "sqlite":
		dialect = sqlite.Open(cli.DSN)
	}
	err := ctx.Run(&Context{
		Debug: cli.LogSQL,
		Config: gorm.Config{
			Logger: logger.Default.LogMode(func() logger.LogLevel {
				if cli.LogSQL {
					return logger.Info
				}
				return logger.Warn
			}()),
			Dialector: dialect,
		},
	})
	ctx.FatalIfErrorf(err)
}

// merge options appends the options to the DSN if they are not already present.
func mergeOptions(dsn, options string) string {
	if options == "" {
		return dsn
	}
	if strings.Contains(dsn, "?") {
		return dsn + "&" + options
	}
	return dsn + "?" + options
}
