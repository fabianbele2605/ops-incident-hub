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

func TestAssignIncidentUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("successful assignment", func(t *testing.T) {
		adminUser, _ := domain.NewUser("admin@test.com", "Admin", domain.RoleAdmin)
		operatorUser, _ := domain.NewUser("operator@test.com", "Operator", domain.RoleOperator)
		testIncident, _ := domain.NewIncident("Test", "Test", domain.SeverityHigh, adminUser.ID)
		
		userRepo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				if id == adminUser.ID {
					return adminUser, nil
				}
				return operatorUser, nil
			},
		}
		
		incidentRepo := &mockIncidentRepository{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
				return testIncident, nil
			},
			updateFunc: func(ctx context.Context, inc *domain.Incident) error {
				return nil
			},
		}
		
		uc := incident.NewAssignIncidentUseCase(incidentRepo, userRepo)
		
		input := incident.AssignIncidentInput{
			IncidentID: testIncident.ID,
			AssignedTo: operatorUser.ID,
			AssignedBy: adminUser.ID,
		}
		
		err := uc.Execute(ctx, input)
		
		require.NoError(t, err)
	})

	t.Run("assigned by user not found", func(t *testing.T) {
		userRepo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return nil, nil
			},
		}
		
		incidentRepo := &mockIncidentRepository{}
		
		uc := incident.NewAssignIncidentUseCase(incidentRepo, userRepo)
		
		input := incident.AssignIncidentInput{
			IncidentID: uuid.New(),
			AssignedTo: uuid.New(),
			AssignedBy: uuid.New(),
		}
		
		err := uc.Execute(ctx, input)
		
		assert.Error(t, err)
		assert.Equal(t, domain.ErrUserNotFound, err)
	})

	t.Run("assigned by user forbidden", func(t *testing.T) {
		viewerUser, _ := domain.NewUser("viewer@test.com", "Viewer", domain.RoleViewer)
		
		userRepo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return viewerUser, nil
			},
		}
		
		incidentRepo := &mockIncidentRepository{}
		
		uc := incident.NewAssignIncidentUseCase(incidentRepo, userRepo)
		
		input := incident.AssignIncidentInput{
			IncidentID: uuid.New(),
			AssignedTo: uuid.New(),
			AssignedBy: viewerUser.ID,
		}
		
		err := uc.Execute(ctx, input)
		
		assert.Error(t, err)
		assert.Equal(t, domain.ErrForbidden, err)
	})

	t.Run("assigned to user not found", func(t *testing.T) {
		adminUser, _ := domain.NewUser("admin@test.com", "Admin", domain.RoleAdmin)
		
		userRepo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				if id == adminUser.ID {
					return adminUser, nil
				}
				return nil, nil
			},
		}
		
		incidentRepo := &mockIncidentRepository{}
		
		uc := incident.NewAssignIncidentUseCase(incidentRepo, userRepo)
		
		input := incident.AssignIncidentInput{
			IncidentID: uuid.New(),
			AssignedTo: uuid.New(),
			AssignedBy: adminUser.ID,
		}
		
		err := uc.Execute(ctx, input)
		
		assert.Error(t, err)
		assert.Equal(t, domain.ErrUserNotFound, err)
	})

	t.Run("incident not found", func(t *testing.T) {
		adminUser, _ := domain.NewUser("admin@test.com", "Admin", domain.RoleAdmin)
		operatorUser, _ := domain.NewUser("operator@test.com", "Operator", domain.RoleOperator)
		
		userRepo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				if id == adminUser.ID {
					return adminUser, nil
				}
				return operatorUser, nil
			},
		}
		
		incidentRepo := &mockIncidentRepository{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
				return nil, nil
			},
		}
		
		uc := incident.NewAssignIncidentUseCase(incidentRepo, userRepo)
		
		input := incident.AssignIncidentInput{
			IncidentID: uuid.New(),
			AssignedTo: operatorUser.ID,
			AssignedBy: adminUser.ID,
		}
		
		err := uc.Execute(ctx, input)
		
		assert.Error(t, err)
		assert.Equal(t, domain.ErrIncidentNotFound, err)
	})

	t.Run("repository update error", func(t *testing.T) {
		adminUser, _ := domain.NewUser("admin@test.com", "Admin", domain.RoleAdmin)
		operatorUser, _ := domain.NewUser("operator@test.com", "Operator", domain.RoleOperator)
		testIncident, _ := domain.NewIncident("Test", "Test", domain.SeverityHigh, adminUser.ID)
		
		userRepo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				if id == adminUser.ID {
					return adminUser, nil
				}
				return operatorUser, nil
			},
		}
		
		repoErr := errors.New("database error")
		incidentRepo := &mockIncidentRepository{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
				return testIncident, nil
			},
			updateFunc: func(ctx context.Context, inc *domain.Incident) error {
				return repoErr
			},
		}
		
		uc := incident.NewAssignIncidentUseCase(incidentRepo, userRepo)
		
		input := incident.AssignIncidentInput{
			IncidentID: testIncident.ID,
			AssignedTo: operatorUser.ID,
			AssignedBy: adminUser.ID,
		}
		
		err := uc.Execute(ctx, input)
		
		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
	})
}
