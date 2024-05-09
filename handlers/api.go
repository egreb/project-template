package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/egreb/boilerplate/internalerrors"
)

type apiError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (a apiError) Error() string {
	return a.Message
}

func (a apiError) Unwrap() error {
	return fmt.Errorf("%d: %s", a.StatusCode, a.Message)
}

func handleApiError(err error) *apiError {
	switch err.(type) {
	case *internalerrors.BadRequest:
		return &apiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad request",
		}
	case *internalerrors.NotFound:
		return &apiError{
			StatusCode: http.StatusNotFound,
			Message:    "Not found",
		}
	case *internalerrors.InternalError:
		return &apiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	case *internalerrors.UnprocessableError:
		return &apiError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Unable to handle invalid data",
		}
	}

	return nil
}

func NewApiError(statusCode int, err error) apiError {
	return apiError{
		StatusCode: statusCode,
		Message:    err.Error(),
	}
}

type apiHandler func(w http.ResponseWriter, r *http.Request) error

func handler(a apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := a(w, r); err != nil {
			if apiErr, ok := err.(apiError); ok {
				writeJSON(w, apiErr.StatusCode, apiErr.Message)
			} else {
				checkErr := handleApiError(err)
				if checkErr != nil {
					writeJSON(w, checkErr.StatusCode, checkErr.Message)
				} else {
					errorResponse := map[string]any{
						"statusCode": http.StatusInternalServerError,
						"message":    "internal server error",
					}
					writeJSON(w, http.StatusInternalServerError, errorResponse)
				}

			}
			slog.Error("HTTP API Error", "err", err.Error(), "path", r.URL.Path)
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
