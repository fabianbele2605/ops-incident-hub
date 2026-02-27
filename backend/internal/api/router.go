package api

import (
	"log/slog"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/api/handler"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/observability/logger"
	"github.com/gorilla/mux"
)

// SetupRouter configura las rutas de la API
func SetupRouter(incidentHandler *handler.IncidentHandler, healthHandler *handler.HealthHandler, appLogger *slog.Logger) *mux.Router {
	router := mux.NewRouter()

	// Middleware global
	router.Use(logger.LoggingMiddleware(appLogger))

	// Health endpoints (antes de /api/v1)
	router.HandleFunc("/health", healthHandler.Health).Methods("GET")
	router.HandleFunc("/health/ready", healthHandler.Ready).Methods("GET")
	router.HandleFunc("/health/live", healthHandler.Live).Methods("GET")

	// API v1
	api := router.PathPrefix("/api/v1").Subrouter()

	// Incidents routes
	api.HandleFunc("/incidents", incidentHandler.Create).Methods("POST")
	api.HandleFunc("/incidents", incidentHandler.List).Methods("GET")
	api.HandleFunc("/incidents/{id}/assign", incidentHandler.Assign).Methods("POST")

	return router
}
