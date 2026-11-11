package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func requestLoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			wrapped := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				status := wrapped.Status()
				if status == 0 {
					status = http.StatusOK
				}
				logger.Info("request", "method", r.Method, "path", r.URL.Path, "status", status, "duration", time.Since(started), "request_id", middleware.GetReqID(r.Context()))
			}()
			next.ServeHTTP(wrapped, r)
		}))
	}
}

func recoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if value := recover(); value != nil {
					logger.Error("request panic", "panic", value, "request_id", middleware.GetReqID(r.Context()))
					writeJSON(w, http.StatusInternalServerError, errorResponse{"internal server error", "INTERNAL_ERROR"})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
