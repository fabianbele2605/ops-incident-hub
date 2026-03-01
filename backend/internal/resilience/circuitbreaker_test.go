package resilience

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_Success(t *testing.T) {
	cb := NewCircuitBreaker(3, 1*time.Second)
	
	err := cb.Execute(func() error {
		return nil
	})
	
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	
	if cb.State() != StateClosed {
		t.Errorf("expected state Closed, got %v", cb.State())
	}
}

func TestCircuitBreaker_OpenAfterFailures(t *testing.T) {
	cb := NewCircuitBreaker(3, 1*time.Second)
	
	// Generar 3 fallos
	for i := 0; i < 3; i++ {
		_ = cb.Execute(func() error {
			return errors.New("error")
		})
	}
	
	if cb.State() != StateOpen {
		t.Errorf("expected state Open, got %v", cb.State())
	}
	
	// Siguiente ejecución debe fallar inmediatamente
	err := cb.Execute(func() error {
		return nil
	})
	
	if err != ErrCircuitOpen {
		t.Errorf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_HalfOpenAfterTimeout(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond)
	
	// Generar fallos para abrir circuito
	_ = cb.Execute(func() error { return errors.New("error") })
	_ = cb.Execute(func() error { return errors.New("error") })
	
	if cb.State() != StateOpen {
		t.Errorf("expected state Open, got %v", cb.State())
	}
	
	// Esperar timeout
	time.Sleep(150 * time.Millisecond)
	
	// Siguiente ejecución debe intentar (HalfOpen)
	err := cb.Execute(func() error {
		return nil
	})
	
	if err != nil {
		t.Errorf("expected no error after timeout, got %v", err)
	}
	
	if cb.State() != StateClosed {
		t.Errorf("expected state Closed after success, got %v", cb.State())
	}
}
