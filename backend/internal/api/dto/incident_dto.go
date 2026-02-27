package dto

import (
	"time"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
)

// CreateIncidentRequest representa la peticion para crear un incidente
type CreateIncidentRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

// AssignIncidentRequest representa la peticion para asignar un incidente
type AssignIncidentRequest struct {
	AssignedTo string `json:"assigned_to"`
}

// IncidentResponse representa la respuesta de un incidente
type IncidentResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Severity    string     `json:"severity"`
	Status      string     `json:"status"`
	AssignedTo  *string    `json:"assigned_to,omitempty"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

// ToIncidentResponse convierte una entidad Incident a DTO
func ToIncidentResponse(incident *domain.Incident) *IncidentResponse {
	response := &IncidentResponse{
		ID:          incident.ID.String(),
		Title:       incident.Title,
		Description: incident.Description,
		Severity:    string(incident.Severity),
		Status:      string(incident.Status),
		CreatedBy:   incident.CreatedBy.String(),
		CreatedAt:   incident.CreatedAt,
		UpdatedAt:   incident.UpdatedAt,
		ResolvedAt:  incident.ResolvedAt,
	}

	if incident.AssignedTo != nil {
		assignedTo := incident.AssignedTo.String()
		response.AssignedTo = &assignedTo
	}

	return response
}

// ToIncidentResponseList convierte una lista de incidentes a DTOs
func ToIncidentResponseList(incidents []*domain.Incident) []*IncidentResponse {
	responses := make([]*IncidentResponse, len(incidents))
	for i, incident := range incidents {
		responses[i] = ToIncidentResponse(incident)
	}
	return responses
}
