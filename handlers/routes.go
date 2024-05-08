package handlers

import (
	"net/http"

	"github.com/egreb/boilerplate/db/repo"
	"github.com/egreb/boilerplate/middleware"
)

func SetupRoutes(router *http.ServeMux, ur *repo.UsersRepository, sr *repo.SessionsRepository) {
	authMiddleware := middleware.Authenticate(sr)

	router.HandleFunc("GET /yo", yohandler)

	me := meHandler(ur)

	router.HandleFunc("POST /api/auth/signin", handler(signinUserHandler(ur, sr)))
	router.HandleFunc("POST /api/auth/register", handler(registerUserHandler(ur)))
	router.HandleFunc("POST /api/auth/signout", handler(signoutUserHandler(sr)))
	router.Handle("GET /api/auth/me", authMiddleware(
		http.HandlerFunc(handler(me)),
	))
	router.Handle("GET /welcome", authMiddleware(
		http.HandlerFunc(handler(func(w http.ResponseWriter, r *http.Request) error {
			return me(w, r)
		})),
	))

}

var yohandler = handler(func(w http.ResponseWriter, r *http.Request) error {
	d := map[string]string{
		"sup": "nothing",
	}

	return writeJSON(w, http.StatusOK, d)
})
