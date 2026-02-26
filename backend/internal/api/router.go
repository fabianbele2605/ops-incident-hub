package api

import (
	"net/http"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/api/handler"
	"github.com/gorilla/mux"
)

// SetupRouter configura las rutas de la API
func SetupRouter(incidentHandler *handler.IncidentHandler) *mux.Router {
	router := mux.NewRouter()

	// API v1
	api := router.PathPrefix("/api/v1").Subrouter()

	// Incidents routes
	api.HandleFunc("/incidents", incidentHandler.Create).Methods("POST")
	api.HandleFunc("/incidents", incidentHandler.List).Methods("GET")
	api.HandleFunc("/incidents/{id}/assign", incidentHandler.Assign).Methods("POST")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return router
}
