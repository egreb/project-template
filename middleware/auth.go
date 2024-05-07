package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/egreb/boilerplate/db/repo"
	"github.com/egreb/boilerplate/errors"
	"github.com/egreb/boilerplate/models"
)

type middlewareContextKey string

func (c middlewareContextKey) String() string {
	return "middleware.middlewareContextKey " + string(c)
}

var (
	contextKeyUser = middlewareContextKey("session-user-id")
)

func getSessionTokenFromContext(ctx context.Context) (string, bool) {
	tokenStr, ok := ctx.Value(contextKeyUser).(string)

	return tokenStr, ok

}

func User(ctx context.Context, ur repo.UsersRepository) (*models.Me, error) {
	tok, ok := getSessionTokenFromContext(ctx)
	if !ok {
		return nil, errors.NotFound{
			Err: fmt.Errorf("user not found"),
		}
	}

	me, err := ur.GetMe(ctx, tok)
	if err != nil {
		return nil, fmt.Errorf("failed getting user from session context: %w", err)
	}

	return me, nil
}

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

			v := context.WithValue(ctx, contextKeyUser, userSession.ID)
			wrappedR := r.WithContext(v)

			next.ServeHTTP(w, wrappedR)
		})
	}
}
