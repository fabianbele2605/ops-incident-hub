package incident

import (
	"context"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/google/uuid"
)

// ListIncidentsInput representa los filtros de busqueda
type ListIncidentsInput struct {
	Status     *domain.Status
	Severity   *domain.Severity
	AssignedTo *uuid.UUID
	Limit      int
	Offset     int
}

// ListIncidentsOutput representa el resultado
type ListIncidentsOutput struct {
	Incidents  []*domain.Incident
	Total      int
}

// ListIncidentsUseCase maneja el listado de incidentes
type ListIncidentsUseCase struct {
	incidentRepo    domain.IncidentRepository
}

// NewListIncidentUseCase crea una nueva instancia
func NewListIncidentUseCase(incidentRepo domain.IncidentRepository) *ListIncidentsUseCase {
	return &ListIncidentsUseCase{
		incidentRepo: incidentRepo,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListIncidentsUseCase) Execute(ctx context.Context, input ListIncidentsInput) (*ListIncidentsOutput, error) {
	// Valores por defecto para paginacion
	if input.Limit <= 0 {
		input.Limit = 20
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	// Crear filtros
	filters := domain.ListFilters{
		Status:       input.Status,
		Severity:     input.Severity,
		AssignedTo:   input.AssignedTo,
		Limit:        input.Limit,
		Offset:       input.Offset,
		SortBy:       "created_at",
		SortOrder:    "desc",
	}

	// Obtener incidentes
	incidents, err := uc.incidentRepo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Contar total
	total, err := uc.incidentRepo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	return &ListIncidentsOutput{
		Incidents: incidents,
		Total:     total,
	}, nil
}