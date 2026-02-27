package domain

import "errors"

var (
	// Incident errors
	ErrInvalidTitle        = errors.New("invalid incident title")
	ErrInvalidDescription  = errors.New("invalid incident description")
	ErrInvalidSeverity     = errors.New("invalid incident severity")
	ErrIncidentNotFound    = errors.New("incident not found")
	ErrIncidentClosed      = errors.New("incident is closed")
	ErrIncidentNotAssigned = errors.New("incident is not assigned")
	ErrIncidentNotResolved = errors.New("incident is not resolved")

	// User errors
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidName       = errors.New("invalid name")
	ErrInvalidRole       = errors.New("invalid role")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	// Generic errors
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrInternalServer = errors.New("internal server error")
)
