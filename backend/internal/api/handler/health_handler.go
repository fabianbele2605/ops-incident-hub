package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse respuesta de health check
type HealthResponse struct {
	Status       string            `json:"status"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
	Timestamp    time.Time         `json:"timestamp"`
}

// Health endpoint básico con verificación de dependencias
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:       "healthy",
		Version:      "1.0.0",
		Dependencies: make(map[string]string),
		Timestamp:    time.Now(),
	}
	
	// Verificar base de datos
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	
	if err := h.db.PingContext(ctx); err != nil {
		response.Status = "unhealthy"
		response.Dependencies["database"] = "down"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	
	response.Dependencies["database"] = "up"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// Live liveness probe (proceso vivo)
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// Ready readiness probe (listo para recibir tráfico)
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// Verificar que todas las dependencias estén listas
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	
	if err := h.db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("Database not ready"))
		return
	}
	
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Ready"))
}
