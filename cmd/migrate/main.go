package main

import (
	"log"
	"os"

	"github.com/egreb/boilerplate/config"
	"github.com/egreb/boilerplate/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := db.GetDB(c.Database)
	if err != nil {
		log.Fatalln(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalln("failed to run migrations: ", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"pgx", driver)
	if err != nil {
		log.Fatalln("failed to run migrations: ", err)
	}

	if len(os.Args) < 2 {
		log.Fatalln("expected 'action' subcommand")
	}

	switch os.Args[1] {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "rollback":
		err = m.Steps(-1)
	default:
		log.Fatalln("no action provided")
	}
	if err != nil {
		log.Fatalln("migrations failed", err)
	}
}
