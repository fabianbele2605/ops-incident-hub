package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/api/dto"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/usecase/incident"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CreateIncidentExecutor interfaz para el caso de uso de creación
type CreateIncidentExecutor interface {
	Execute(ctx context.Context, input incident.CreateIncidentInput) (*domain.Incident, error)
}

// AssignIncidentExecutor interfaz para el caso de uso de asignación
type AssignIncidentExecutor interface {
	Execute(ctx context.Context, input incident.AssignIncidentInput) error
}

// ListIncidentsExecutor interfaz para el caso de uso de listado
type ListIncidentsExecutor interface {
	Execute(ctx context.Context, input incident.ListIncidentsInput) (*incident.ListIncidentsOutput, error)
}

// IncidentHandler maneja las peticiones HTTP de incidentes
type IncidentHandler struct {
	createUseCase CreateIncidentExecutor
	assignUseCase AssignIncidentExecutor
	listUseCase   ListIncidentsExecutor
}

// NewIncidentHandler crea una nueva instancia
func NewIncidentHandler(
	createUseCase CreateIncidentExecutor,
	assignUseCase AssignIncidentExecutor,
	listUseCase ListIncidentsExecutor,
) *IncidentHandler {
	return &IncidentHandler{
		createUseCase: createUseCase,
		assignUseCase: assignUseCase,
		listUseCase:   listUseCase,
	}
}

// Create maneja POST /api/v1/incidents
func (h *IncidentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateIncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Obtener usuario del contexto (autenticación)
	createdBy := uuid.New() // Temporal

	input := incident.CreateIncidentInput{
		Title:       req.Title,
		Description: req.Description,
		Severity:    domain.Severity(req.Severity),
		CreatedBy:   createdBy,
	}

	createdIncident, err := h.createUseCase.Execute(r.Context(), input)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, dto.ToIncidentResponse(createdIncident))
}

// Assign maneja POST /api/v1/incidents/:id/assign
func (h *IncidentHandler) Assign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	incidentID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid incident ID")
		return
	}

	var req dto.AssignIncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	assignedTo, err := uuid.Parse(req.AssignedTo)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid assigned_to ID")
		return
	}

	// TODO: Obtener usuario del contexto
	assignedBy := uuid.New() // Temporal

	input := incident.AssignIncidentInput{
		IncidentID: incidentID,
		AssignedTo: assignedTo,
		AssignedBy: assignedBy,
	}

	if err := h.assignUseCase.Execute(r.Context(), input); err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, dto.SuccessResponse{
		Message: "Incident assigned successfully",
	})
}

// List maneja GET /api/v1/incidents
func (h *IncidentHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: Parsear query params para filtros
	input := incident.ListIncidentsInput{
		Limit:  20,
		Offset: 0,
	}

	output, err := h.listUseCase.Execute(r.Context(), input)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	totalPages := (output.Total + input.Limit - 1) / input.Limit

	response := dto.PaginatedResponse{
		Data:       dto.ToIncidentResponseList(output.Incidents),
		Total:      output.Total,
		Limit:      input.Limit,
		Offset:     input.Offset,
		TotalPages: totalPages,
	}

	respondJSON(w, http.StatusOK, response)
}

// respondJSON envía una respuesta JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// respondError envía una respuesta de error
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, dto.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}

// handleUseCaseError maneja errores de casos de uso
func handleUseCaseError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrInvalidTitle, domain.ErrInvalidDescription, domain.ErrInvalidSeverity:
		respondError(w, http.StatusBadRequest, err.Error())
	case domain.ErrIncidentNotFound, domain.ErrUserNotFound:
		respondError(w, http.StatusNotFound, err.Error())
	case domain.ErrForbidden:
		respondError(w, http.StatusForbidden, err.Error())
	case domain.ErrIncidentClosed, domain.ErrIncidentNotAssigned, domain.ErrIncidentNotResolved:
		respondError(w, http.StatusConflict, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, "Internal server error")
	}
}
