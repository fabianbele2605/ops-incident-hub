package metrics

import (
	"net/http"
	"strconv"
	"time"
)

// MetricsMiddleware middleware para registrar métricas HTTP
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrapper para capturar status code
		wrapped := &metricsResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Ejecutar handler
		next.ServeHTTP(wrapped, r)

		// Registrar métricas
		duration := time.Since(start)
		endpoint := normalizeEndpoint(r.URL.Path)
		status := strconv.Itoa(wrapped.statusCode)

		RecordHTTPRequest(r.Method, endpoint, status, duration)
	})
}

// metricsResponseWriter wrapper para capturar status code
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (mrw *metricsResponseWriter) WriteHeader(code int) {
	mrw.statusCode = code
	mrw.ResponseWriter.WriteHeader(code)
}

// normalizeEndpoint normaliza endpoints para evitar alta cardinalidad
func normalizeEndpoint(path string) string {
	// Normalizar paths con IDs
	// /api/v1/incidents/123 -> /api/v1/incidents/{id}
	// /api/v1/incidents/123/assign -> /api/v1/incidents/{id}/assign
	
	switch {
	case path == "/health":
		return "/health"
	case path == "/health/ready":
		return "/health/ready"
	case path == "/health/live":
		return "/health/live"
	case path == "/metrics":
		return "/metrics"
	case path == "/api/v1/incidents":
		return "/api/v1/incidents"
	default:
		// Si contiene /incidents/ seguido de algo, normalizar
		if len(path) > 20 && path[:20] == "/api/v1/incidents/" {
			if len(path) > 57 && path[len(path)-7:] == "/assign" {
				return "/api/v1/incidents/{id}/assign"
			}
			return "/api/v1/incidents/{id}"
		}
		return path
	}
}
