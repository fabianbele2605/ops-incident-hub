package incident_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/usecase/incident"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockIncidentRepository es un mock del repositorio de incidentes
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

// mockUserRepository es un mock del repositorio de usuarios
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

func TestCreateIncidentUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		adminUser, _ := domain.NewUser("admin@test.com", "Admin", domain.RoleAdmin)
		
		userRepo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return adminUser, nil
			},
		}
		
		incidentRepo := &mockIncidentRepository{
			createFunc: func(ctx context.Context, inc *domain.Incident) error {
				return nil
			},
		}
		
		uc := incident.NewCreateIncidentUseCase(incidentRepo, userRepo)
		
		input := incident.CreateIncidentInput{
			Title:       "Test Incident",
			Description: "Test Description",
			Severity:    domain.SeverityHigh,
			CreatedBy:   adminUser.ID,
		}
		
		result, err := uc.Execute(ctx, input)
		
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Incident", result.Title)
		assert.Equal(t, domain.StatusOpen, result.Status)
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return nil, nil
			},
		}
		
		incidentRepo := &mockIncidentRepository{}
		
		uc := incident.NewCreateIncidentUseCase(incidentRepo, userRepo)
		
		input := incident.CreateIncidentInput{
			Title:       "Test",
			Description: "Test",
			Severity:    domain.SeverityHigh,
			CreatedBy:   uuid.New(),
		}
		
		result, err := uc.Execute(ctx, input)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrUserNotFound, err)
	})

	t.Run("user cannot manage incidents", func(t *testing.T) {
		viewerUser, _ := domain.NewUser("viewer@test.com", "Viewer", domain.RoleViewer)
		
		userRepo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return viewerUser, nil
			},
		}
		
		incidentRepo := &mockIncidentRepository{}
		
		uc := incident.NewCreateIncidentUseCase(incidentRepo, userRepo)
		
		input := incident.CreateIncidentInput{
			Title:       "Test",
			Description: "Test",
			Severity:    domain.SeverityHigh,
			CreatedBy:   viewerUser.ID,
		}
		
		result, err := uc.Execute(ctx, input)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrForbidden, err)
	})

	t.Run("invalid incident data", func(t *testing.T) {
		adminUser, _ := domain.NewUser("admin@test.com", "Admin", domain.RoleAdmin)
		
		userRepo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return adminUser, nil
			},
		}
		
		incidentRepo := &mockIncidentRepository{}
		
		uc := incident.NewCreateIncidentUseCase(incidentRepo, userRepo)
		
		input := incident.CreateIncidentInput{
			Title:       "",
			Description: "Test",
			Severity:    domain.SeverityHigh,
			CreatedBy:   adminUser.ID,
		}
		
		result, err := uc.Execute(ctx, input)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrInvalidTitle, err)
	})

	t.Run("repository error", func(t *testing.T) {
		adminUser, _ := domain.NewUser("admin@test.com", "Admin", domain.RoleAdmin)
		
		userRepo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return adminUser, nil
			},
		}
		
		repoErr := errors.New("database error")
		incidentRepo := &mockIncidentRepository{
			createFunc: func(ctx context.Context, inc *domain.Incident) error {
				return repoErr
			},
		}
		
		uc := incident.NewCreateIncidentUseCase(incidentRepo, userRepo)
		
		input := incident.CreateIncidentInput{
			Title:       "Test",
			Description: "Test",
			Severity:    domain.SeverityHigh,
			CreatedBy:   adminUser.ID,
		}
		
		result, err := uc.Execute(ctx, input)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, repoErr, err)
	})
}
