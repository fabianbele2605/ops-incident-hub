package logger

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// LoggingMiddleware middleware para logging de requests HTTP
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Generar request ID
			requestID := fmt.Sprintf("req-%s", uuid.New().String()[:8])

			// Crear logger con request_id
			reqLogger := logger.With(
				"request_id", requestID,
				"method", r.Method,
				"path", r.URL.Path,
				"ip", r.RemoteAddr,
			)

			// Agregar logger al contexto
			ctx := WithContext(r.Context(), reqLogger)
			r = r.WithContext(ctx)

			// Agregar request_id al response header
			w.Header().Set("X-Request-ID", requestID)

			// Log inicio de request
			reqLogger.Info("http request started")

			// Wrapper para capturar status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Recuperar de panics
			defer func() {
				if err := recover(); err != nil {
					reqLogger.Error("panic recovered",
						"error", err,
						"duration_ms", time.Since(start).Milliseconds(),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			// Ejecutar handler
			next.ServeHTTP(wrapped, r)

			// Log fin de request
			reqLogger.Info("http request completed",
				"status", wrapped.statusCode,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}

// responseWriter wrapper para capturar status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
