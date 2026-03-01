package resilience

import (
	"context"
	"errors"
	"math"
	"time"
)

// RetryConfig configuración de retry
type RetryConfig struct {
	MaxAttempts int
	InitialWait time.Duration
	MaxWait     time.Duration
	Multiplier  float64
}

// DefaultRetryConfig configuración por defecto
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		InitialWait: 100 * time.Millisecond,
		MaxWait:     5 * time.Second,
		Multiplier:  2.0,
	}
}

// IsRetryable función que determina si un error es reintentable
type IsRetryable func(error) bool

// Retry ejecuta una función con reintentos y backoff exponencial
func Retry(ctx context.Context, cfg RetryConfig, isRetryable IsRetryable, fn func() error) error {
	var lastErr error
	
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		// Ejecutar función
		err := fn()
		
		// Si no hay error, retornar éxito
		if err == nil {
			return nil
		}
		
		lastErr = err
		
		// Si no es reintentable, retornar error inmediatamente
		if isRetryable != nil && !isRetryable(err) {
			return err
		}
		
		// Si es el último intento, no esperar
		if attempt == cfg.MaxAttempts-1 {
			break
		}
		
		// Calcular tiempo de espera con backoff exponencial
		wait := calculateBackoff(attempt, cfg)
		
		// Esperar con respeto al contexto
		select {
		case <-time.After(wait):
			// Continuar con siguiente intento
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	
	return lastErr
}

// calculateBackoff calcula el tiempo de espera con backoff exponencial
func calculateBackoff(attempt int, cfg RetryConfig) time.Duration {
	wait := float64(cfg.InitialWait) * math.Pow(cfg.Multiplier, float64(attempt))
	
	if wait > float64(cfg.MaxWait) {
		wait = float64(cfg.MaxWait)
	}
	
	return time.Duration(wait)
}

// IsTemporaryError verifica si un error es temporal
func IsTemporaryError(err error) bool {
	// Errores temporales comunes
	var tempErr interface{ Temporary() bool }
	if errors.As(err, &tempErr) {
		return tempErr.Temporary()
	}
	
	// Agregar más casos según necesidad
	return false
}
