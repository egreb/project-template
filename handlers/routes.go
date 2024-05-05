package handlers

import (
	"net/http"

	"github.com/egreb/boilerplate/db/repo"
	"github.com/egreb/boilerplate/middleware"
)

func SetupRoutes(router *http.ServeMux, ur repo.UsersRepository, sr repo.SessionsRepository) {
	authMiddleware := middleware.Authenticate(sr)

	router.HandleFunc("GET /yo", yohandler)

	router.HandleFunc("POST /api/auth/signin", handler(signinUserHandler(ur, sr)))
	router.HandleFunc("POST /api/auth/register", handler(registerUserHandler(ur)))
	router.Handle("GET /welcome", authMiddleware(
		http.HandlerFunc(handler(func(w http.ResponseWriter, r *http.Request) error {
			return writeJSON(w, http.StatusOK, "Welcome")
		})),
	))
}

var yohandler = handler(func(w http.ResponseWriter, r *http.Request) error {
	d := map[string]string{
		"sup": "nothing",
	}

	return writeJSON(w, http.StatusOK, d)
})
