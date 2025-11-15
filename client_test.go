package togomq

import (
	"context"
	"testing"
)

func TestCountMessages_Validation(t *testing.T) {
	// Create a client with default config (won't actually connect)
	cfg := NewConfig(WithToken("test-token"))
	client := &Client{
		config: cfg,
		logger: NewLogger(LogLevelNone),
	}

	// Test empty topic validation
	t.Run("empty topic", func(t *testing.T) {
		_, err := client.CountMessages(context.Background(), "")
		if err == nil {
			t.Error("Expected error for empty topic, got nil")
			return
		}
		tmqErr, ok := err.(*TogoMQError)
		if !ok {
			t.Errorf("Expected TogoMQError, got %T", err)
			return
		}
		if tmqErr.Code != ErrCodeValidation {
			t.Errorf("Expected error code %s, got %s", ErrCodeValidation, tmqErr.Code)
		}
	})
}
