package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/egreb/boilerplate/auth"
	"github.com/egreb/boilerplate/db/repo"
	"github.com/egreb/boilerplate/internalerrors"
	"github.com/egreb/boilerplate/middleware"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func registerUserHandler(ur *repo.UsersRepository) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var creds credentials

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			return &internalerrors.UnprocessableError{}
		}

		salt, err := auth.GenerateRandomSalt(24)
		if err != nil {
			return err
		}
		hashedPassword := auth.HashPassword(creds.Password, salt)

		err = ur.CreateUser(r.Context(), creds.Username, hashedPassword, salt)
		if err != nil {
			return err
		}

		return writeJSON(w, http.StatusCreated, "success")
	}
}

func signinUserHandler(ur *repo.UsersRepository, sr *repo.SessionsRepository) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var creds credentials

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			return &internalerrors.UnprocessableError{}
		}

		id, hashedPassword, salt, err := ur.GetUserCredentialsByUsername(r.Context(), creds.Username)
		if err != nil {
			return err
		}

		ok := auth.ComparePasswords(hashedPassword, creds.Password, salt)
		if !ok {
			return NewApiError(http.StatusBadRequest, fmt.Errorf("failed to login, please try again"))
		}

		session, err := sr.CreateSession(r.Context(), id)
		if err != nil {
			return err
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    session.ID,
			Expires:  session.Expires,
			HttpOnly: true,
			SameSite: http.SameSiteDefaultMode,
		})

		return writeJSON(w, http.StatusOK, "ok")
	}
}

func signoutUserHandler(sr *repo.SessionsRepository) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		tok, ok := middleware.SessionToken(r.Context())
		if !ok {
			return NewApiError(http.StatusNotFound, fmt.Errorf("session not found"))
		}

		err := sr.Delete(r.Context(), tok)
		if err != nil {
			return err
		}

		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				return writeJSON(w, http.StatusNotFound, "This makes no sense")
			}
			return err
		}

		c.Expires = time.Now()
		http.SetCookie(w, c)

		return writeJSON(w, http.StatusOK, "ok")
	}
}

func meHandler(ur *repo.UsersRepository) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		me, err := middleware.User(r.Context(), ur)
		if err != nil {
			return err
		}

		return writeJSON(w, http.StatusOK, me)
	}
}
