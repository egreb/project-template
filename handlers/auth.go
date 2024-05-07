package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/egreb/boilerplate/auth"
	"github.com/egreb/boilerplate/db/repo"
	"github.com/egreb/boilerplate/errors"
	"github.com/egreb/boilerplate/middleware"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func registerUserHandler(ur repo.UsersRepository) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var creds credentials

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			return errors.BadRequest{
				Err: fmt.Errorf("error parsing register request: %w", err),
			}
		}

		salt, err := auth.GenerateRandomSalt(24)
		if err != nil {
			// TODO: Write handleerr function for this to do logging
			return writeJSON(w, http.StatusInternalServerError, "something went wrong")
		}
		hashedPassword := auth.HashPassword(creds.Password, salt)

		err = ur.CreateUser(r.Context(), creds.Username, hashedPassword, salt)
		if err != nil {
			return errors.InternalError{
				Err: fmt.Errorf("error storing user: %w", err),
			}
		}

		return writeJSON(w, http.StatusCreated, "success")
	}
}

func signinUserHandler(ur repo.UsersRepository, sr repo.SessionsRepository) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var creds credentials

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			return errors.BadRequest{
				Err: fmt.Errorf("error parsing register request: %w", err),
			}
		}

		id, hashedPassword, salt, err := ur.GetUserCredentialsByUsername(r.Context(), creds.Username)
		if err != nil {
			return fmt.Errorf("unable to signin: %w", err)
		}

		ok := auth.ComparePasswords(hashedPassword, creds.Password, salt)
		if !ok {
			return writeJSON(w, http.StatusBadRequest, "Username or password is wrong")
		}

		session, err := sr.Create(r.Context(), id)
		if err != nil {
			return writeJSON(w, http.StatusInternalServerError, "Something went wrong")
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

func meHandler(ur repo.UsersRepository) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		me, err := middleware.User(r.Context(), ur)
		if err != nil {
			return fmt.Errorf("could not get user from session: %w", err)
		}

		return writeJSON(w, http.StatusOK, me)
	}
}
