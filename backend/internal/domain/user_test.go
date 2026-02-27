package domain_test

import (
	"testing"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		user, err := domain.NewUser("test@example.com", "Test User", domain.RoleAdmin)

		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEqual(t, uuid.Nil, user.ID)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, domain.RoleAdmin, user.Role)
		assert.NotZero(t, user.CreatedAt)
	})

	t.Run("empty email", func(t *testing.T) {
		user, err := domain.NewUser("", "Test User", domain.RoleAdmin)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, domain.ErrInvalidEmail, err)
	})

	t.Run("empty name", func(t *testing.T) {
		user, err := domain.NewUser("test@example.com", "", domain.RoleAdmin)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, domain.ErrInvalidName, err)
	})

	t.Run("invalid role", func(t *testing.T) {
		user, err := domain.NewUser("test@example.com", "Test User", "invalid")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, domain.ErrInvalidRole, err)
	})
}

func TestUser_CanManageIncidents(t *testing.T) {
	t.Run("admin can manage", func(t *testing.T) {
		user, _ := domain.NewUser("admin@example.com", "Admin", domain.RoleAdmin)

		assert.True(t, user.CanManageIncidents())
	})

	t.Run("operator can manage", func(t *testing.T) {
		user, _ := domain.NewUser("operator@example.com", "Operator", domain.RoleOperator)

		assert.True(t, user.CanManageIncidents())
	})

	t.Run("viewer cannot manage", func(t *testing.T) {
		user, _ := domain.NewUser("viewer@example.com", "Viewer", domain.RoleViewer)

		assert.False(t, user.CanManageIncidents())
	})
}
