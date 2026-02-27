package handler_test

import (
	"context"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/usecase/incident"
)

// MockCreateIncidentUseCase mock del caso de uso de creación
type MockCreateIncidentUseCase struct {
	ExecuteFunc func(ctx context.Context, input incident.CreateIncidentInput) (*domain.Incident, error)
}

func (m *MockCreateIncidentUseCase) Execute(ctx context.Context, input incident.CreateIncidentInput) (*domain.Incident, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

// MockAssignIncidentUseCase mock del caso de uso de asignación
type MockAssignIncidentUseCase struct {
	ExecuteFunc func(ctx context.Context, input incident.AssignIncidentInput) error
}

func (m *MockAssignIncidentUseCase) Execute(ctx context.Context, input incident.AssignIncidentInput) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil
}

// MockListIncidentsUseCase mock del caso de uso de listado
type MockListIncidentsUseCase struct {
	ExecuteFunc func(ctx context.Context, input incident.ListIncidentsInput) (*incident.ListIncidentsOutput, error)
}

func (m *MockListIncidentsUseCase) Execute(ctx context.Context, input incident.ListIncidentsInput) (*incident.ListIncidentsOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}
