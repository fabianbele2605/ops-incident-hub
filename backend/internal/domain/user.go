package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role representa el rol de un usuario
type Role string

const (
	RoleAdmin       Role = "admin"
	RoleOperator    Role = "operator"
	RoleViewer      Role = "viewer"
)

// User representa un usuario del sistema
type User struct {
	ID              uuid.UUID
	Email           string
	Name            string
	Role            Role
	CreatedAt       time.Time
}

// NewUser crea un nuevo usuario
func NewUser(email, name string, role Role) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}
	if name == "" {
		return nil, ErrInvalidName
	}
	if !isValidRole(role) {
		return nil, ErrInvalidRole
	}

	return &User{
		ID:             uuid.New(),
		Email:          email,
		Name:           name,
		Role:           role,
		CreatedAt:      time.Now(),
	}, nil
}

// CanViewIncidents verifica si el usuario puede ver incidentes
func (u *User) CanManageIncidents() bool {
	return u.Role == RoleAdmin || u.Role == RoleOperator
}

// CanViewIncidents verifica si el usuario puede ver incidentes
func (u *User) CanViewIncidents() bool {
	return true
}

// IsAdmin verifica si el usuario es administrador
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// isValidRole valida si el rol es valido
func isValidRole(r Role) bool {
	switch r {
	case RoleAdmin, RoleOperator, RoleViewer:
		return true
	default:
		return false
	}
}