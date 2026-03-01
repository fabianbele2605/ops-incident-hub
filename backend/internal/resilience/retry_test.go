package resilience

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetry_Success(t *testing.T) {
	cfg := DefaultRetryConfig()
	ctx := context.Background()
	
	attempts := 0
	err := Retry(ctx, cfg, nil, func() error {
		attempts++
		return nil
	})
	
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestRetry_SuccessAfterFailures(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts: 3,
		InitialWait: 10 * time.Millisecond,
		MaxWait:     100 * time.Millisecond,
		Multiplier:  2.0,
	}
	ctx := context.Background()
	
	attempts := 0
	err := Retry(ctx, cfg, nil, func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	})
	
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetry_MaxAttemptsExceeded(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts: 3,
		InitialWait: 10 * time.Millisecond,
		MaxWait:     100 * time.Millisecond,
		Multiplier:  2.0,
	}
	ctx := context.Background()
	
	attempts := 0
	err := Retry(ctx, cfg, nil, func() error {
		attempts++
		return errors.New("persistent error")
	})
	
	if err == nil {
		t.Error("expected error, got nil")
	}
	
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetry_ContextCanceled(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts: 5,
		InitialWait: 100 * time.Millisecond,
		MaxWait:     1 * time.Second,
		Multiplier:  2.0,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	
	attempts := 0
	err := Retry(ctx, cfg, nil, func() error {
		attempts++
		return errors.New("error")
	})
	
	if err != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
}

func TestRetry_NonRetryableError(t *testing.T) {
	cfg := DefaultRetryConfig()
	ctx := context.Background()
	
	nonRetryableErr := errors.New("non-retryable error")
	
	attempts := 0
	err := Retry(ctx, cfg, func(err error) bool {
		return err != nonRetryableErr
	}, func() error {
		attempts++
		return nonRetryableErr
	})
	
	if err != nonRetryableErr {
		t.Errorf("expected non-retryable error, got %v", err)
	}
	
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}
