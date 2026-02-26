package domain

import (
	"time"

	"github.com/google/uuid"
)

// Severity representa la severidad de un incidente
type Severity string

const (
	SeverityCritical    Severity = "critical"
	SeverityHigh        Severity = "high"
	SeverityMedium      Severity = "medium"
	SeverityLow         Severity = "low"
)

// Status representa el estado de un incidente
type Status string

const (
	StatusOpen      Status = "open"
	StatusAssigned  Status = "assigned"
	StatusInProgress Status = "in_progress"
	StatusResolved   Status = "resolved"
	StatusClosed     Status = "closed"
)

// Incident representa un incidente operativo
type Incident struct {
	ID             uuid.UUID
	Title          string
	Description    string
	Severity       Severity
	Status         Status
	AssignedTo     *uuid.UUID
	CreatedBy      uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ResolvedAt     *time.Time
	Metadata       map[string]interface{}
}

// NewIncident crea un nuevo incidente
func NewIncident(title, description string, severity Severity, createdBy uuid.UUID) (*Incident, error) {
	if title == "" {
		return nil, ErrInvalidTitle
	}
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if !isValidSeverity(severity) {
		return nil, ErrInvalidSeverity
	}

	now := time.Now()
	return &Incident{
		ID:              uuid.New(),
		Title:           title,
		Description:     description,
		Severity:        severity,
		Status:          StatusOpen,
		CreatedBy:       createdBy,
		CreatedAt:       now,
		UpdatedAt:       now,
		Metadata:        make(map[string]interface{}),
	}, nil
}

// Assign asigna el incidente a un usuario
func (i *Incident) Assign(userID uuid.UUID) error {
	if i.Status == StatusClosed {
		return ErrIncidentClosed
	}
	i.AssignedTo = &userID
	i.Status = StatusAssigned
	i.UpdatedAt = time.Now()
	return nil
}

// StartProgress marca el incidente como en progreso
func (i *Incident) StartProgress() error {
	if i.Status == StatusClosed {
		return ErrIncidentClosed
	}
	if i.AssignedTo == nil {
		return ErrIncidentNotAssigned
	}
	i.Status = StatusInProgress
	i.UpdatedAt = time.Now()
	return nil
}

// Resolve marca el incidente como resuelto
func (i *Incident) Resolve() error {
	if i.Status == StatusClosed {
		return ErrIncidentClosed
	}
	now := time.Now()
	i.Status = StatusResolved
	i.ResolvedAt = &now
	i.UpdatedAt = now
	return nil
}

// Close cierra el incidente
func (i *Incident) Close() error {
	if i.Status != StatusResolved {
		return ErrIncidentNotResolved
	}
	i.Status = StatusClosed
	i.UpdatedAt = time.Now()
	return nil
}

// IsOpen verifica si el incidente esta abierto
func (i *Incident) IsOpen() bool {
	return i.Status == StatusOpen || i.Status == StatusAssigned || i.Status == StatusInProgress
}

// isValidSeverity valida si la severidad es valida
func isValidSeverity(s Severity) bool {
	switch s {
	case SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow:
		return true
	default:
		return false
	}
}