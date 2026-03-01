package resilience

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

// State representa el estado del circuit breaker
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker implementa el patrón circuit breaker
type CircuitBreaker struct {
	maxFailures  int
	resetTimeout time.Duration
	
	mu           sync.RWMutex
	state        State
	failures     int
	lastFailTime time.Time
}

// NewCircuitBreaker crea un nuevo circuit breaker
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
	}
}

// Execute ejecuta una función con protección de circuit breaker
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.canExecute() {
		return ErrCircuitOpen
	}
	
	err := fn()
	cb.recordResult(err)
	
	return err
}

// canExecute verifica si se puede ejecutar la función
func (cb *CircuitBreaker) canExecute() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailTime) > cb.resetTimeout {
			cb.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

// recordResult registra el resultado de la ejecución
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	if err != nil {
		cb.failures++
		cb.lastFailTime = time.Now()
		
		if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
		}
	} else {
		cb.failures = 0
		cb.state = StateClosed
	}
}

// State retorna el estado actual
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}
