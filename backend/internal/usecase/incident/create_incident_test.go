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
