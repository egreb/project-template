package main

import (
	"context"
	"log"
	"net/http"

	"github.com/egreb/boilerplate/config"
	"github.com/egreb/boilerplate/db"
	"github.com/egreb/boilerplate/db/repo"
	"github.com/egreb/boilerplate/handlers"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	db, err := db.Connect(ctx, c.Database)
	if err != nil {
		log.Fatalln(err)
	}

	usersRepository := repo.NewUsersRepository(db)
	sessionsRepository := repo.NewSessionsRespository(db)

	r := http.NewServeMux()
	handlers.SetupRoutes(r, *usersRepository, *sessionsRepository)

	log.Println("server listening at port:", c.ServerPort)
	http.ListenAndServe(":"+c.ServerPort, r)

}
