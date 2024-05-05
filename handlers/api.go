package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/egreb/boilerplate/errors"
)

type apiHandler func(w http.ResponseWriter, r *http.Request) error

func handler(a apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := a(w, r); err != nil {
			slog.Error("api error", err)
			status := http.StatusInternalServerError

			switch err.(type) {
			case errors.BadRequest:
				status = http.StatusBadRequest
			case errors.NotFound:
				status = http.StatusNotFound
			case errors.InternalError:
				status = http.StatusInternalServerError
			}
			writeJSON(w, status, err.Error())
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("writing json failed: %w", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(b)
	return nil
}
