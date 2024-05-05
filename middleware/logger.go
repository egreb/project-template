package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type wrapperWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrapperWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logger(next http.Handler) http.Handler {
	// Set up the default logger to be JsonFormatted
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrapperWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		slog.Info("Request",
			"statusCode", fmt.Sprintf("%d", wrapped.statusCode),
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"responseTime", fmt.Sprintf("%v µs", time.Since(start).Microseconds()),
		)
	})
}
