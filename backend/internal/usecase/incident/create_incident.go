package incident

import (
	"context"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/observability/metrics"
	"github.com/google/uuid"
)

// CreateIncidentInput representa los datos de entrada para crear un incidente
type CreateIncidentInput struct {
	Title       string
	Description string
	Severity    domain.Severity
	CreatedBy   uuid.UUID
}

// CreateIncidentUseCase maneja la creacion de incidentes
type CreateIncidentUseCase struct {
	incidentRepo domain.IncidentRepository
	userRepo     domain.UserRepository
}

// NewCreateIncidentUseCase crea una nueva instancia del caso de uso
func NewCreateIncidentUseCase(
	incidentRepo domain.IncidentRepository,
	userRepo domain.UserRepository,
) *CreateIncidentUseCase {
	return &CreateIncidentUseCase{
		incidentRepo: incidentRepo,
		userRepo:     userRepo,
	}
}

// Execute ejecuta el caso de uso de crear incidente
func (uc *CreateIncidentUseCase) Execute(ctx context.Context, input CreateIncidentInput) (*domain.Incident, error) {
	// Validar que el usuario existe
	user, err := uc.userRepo.GetByID(ctx, input.CreatedBy)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	// Validar que el usuario puede crear incidentes
	if !user.CanManageIncidents() {
		return nil, domain.ErrForbidden
	}

	// Crear la entidad de incidente
	incident, err := domain.NewIncident(
		input.Title,
		input.Description,
		input.Severity,
		input.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	// Guardar en el repositorio
	if err := uc.incidentRepo.Create(ctx, incident); err != nil {
		return nil, err
	}

	// Registrar métrica
	metrics.RecordIncidentCreated(string(incident.Severity))

	return incident, nil
}
