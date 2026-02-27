package incident_test

import (
	"context"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/google/uuid"
)

type mockIncidentRepository struct {
	createFunc  func(ctx context.Context, inc *domain.Incident) error
	getByIDFunc func(ctx context.Context, id uuid.UUID) (*domain.Incident, error)
	listFunc    func(ctx context.Context, filters domain.ListFilters) ([]*domain.Incident, error)
	updateFunc  func(ctx context.Context, inc *domain.Incident) error
	deleteFunc  func(ctx context.Context, id uuid.UUID) error
	countFunc   func(ctx context.Context, filters domain.ListFilters) (int, error)
}

func (m *mockIncidentRepository) Create(ctx context.Context, inc *domain.Incident) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, inc)
	}
	return nil
}

func (m *mockIncidentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockIncidentRepository) List(ctx context.Context, filters domain.ListFilters) ([]*domain.Incident, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, filters)
	}
	return nil, nil
}

func (m *mockIncidentRepository) Update(ctx context.Context, inc *domain.Incident) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, inc)
	}
	return nil
}

func (m *mockIncidentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockIncidentRepository) Count(ctx context.Context, filters domain.ListFilters) (int, error) {
	if m.countFunc != nil {
		return m.countFunc(ctx, filters)
	}
	return 0, nil
}

type mockUserRepository struct {
	createFunc     func(ctx context.Context, user *domain.User) error
	getByIDFunc    func(ctx context.Context, id uuid.UUID) (*domain.User, error)
	getByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
	listFunc       func(ctx context.Context) ([]*domain.User, error)
	updateFunc     func(ctx context.Context, user *domain.User) error
	deleteFunc     func(ctx context.Context, id uuid.UUID) error
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *mockUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx)
	}
	return nil, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, user)
	}
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}
