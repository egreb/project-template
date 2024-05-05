package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/egreb/boilerplate/db/repo"
)

func Authenticate(sr repo.SessionsRepository) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			c, err := r.Cookie("session_token")
			if err != nil {
				if err == http.ErrNoCookie {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				w.WriteHeader(http.StatusBadRequest)
				return
			}

			sessionToken := c.Value

			userSession, err := sr.Get(ctx, sessionToken)
			if err != nil {
				slog.Warn("auth", "middleware", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if userSession.Expired() {
				err = sr.Delete(ctx, sessionToken)
				if err != nil {
					slog.Warn("auth", "middleware", "remove", "session", err)
				}
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			v := context.WithValue(ctx, "user", userSession.ID)
			wrappedR := r.WithContext(v)

			next.ServeHTTP(w, wrappedR)
		})
	}
}
