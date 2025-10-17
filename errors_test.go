package togomq

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewError(t *testing.T) {
	code := ErrCodeConnection
	message := "connection failed"
	underlyingErr := errors.New("underlying error")

	err := NewError(code, message, underlyingErr)

	if err.Code != code {
		t.Errorf("Expected code '%s', got '%s'", code, err.Code)
	}
	if err.Message != message {
		t.Errorf("Expected message '%s', got '%s'", message, err.Message)
	}
	if err.Err != underlyingErr {
		t.Errorf("Expected underlying error to be set")
	}
}

func TestTogoMQError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *TogoMQError
		expected string
	}{
		{
			name: "with underlying error",
			err: &TogoMQError{
				Code:    ErrCodeConnection,
				Message: "connection failed",
				Err:     errors.New("timeout"),
			},
			expected: "[CONNECTION_ERROR] connection failed: timeout",
		},
		{
			name: "without underlying error",
			err: &TogoMQError{
				Code:    ErrCodeValidation,
				Message: "invalid input",
				Err:     nil,
			},
			expected: "[VALIDATION_ERROR] invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Expected error string '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestTogoMQError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := NewError(ErrCodeConnection, "connection failed", underlyingErr)

	unwrapped := err.Unwrap()
	if unwrapped != underlyingErr {
		t.Errorf("Expected unwrapped error to be the underlying error")
	}
}

func TestWrapGRPCError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		context      string
		expectedNil  bool
		expectedCode string
	}{
		{
			name:        "nil error",
			err:         nil,
			context:     "test context",
			expectedNil: true,
		},
		{
			name:         "unauthenticated error",
			err:          status.Error(codes.Unauthenticated, "invalid credentials"),
			context:      "authentication failed",
			expectedNil:  false,
			expectedCode: ErrCodeAuth,
		},
		{
			name:         "invalid argument error",
			err:          status.Error(codes.InvalidArgument, "bad request"),
			context:      "validation failed",
			expectedNil:  false,
			expectedCode: ErrCodeValidation,
		},
		{
			name:         "unavailable error",
			err:          status.Error(codes.Unavailable, "service unavailable"),
			context:      "connection failed",
			expectedNil:  false,
			expectedCode: ErrCodeConnection,
		},
		{
			name:         "other gRPC error",
			err:          status.Error(codes.Internal, "internal error"),
			context:      "stream error",
			expectedNil:  false,
			expectedCode: ErrCodeStream,
		},
		{
			name:         "non-gRPC error",
			err:          errors.New("some error"),
			context:      "general error",
			expectedNil:  false,
			expectedCode: ErrCodeStream,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapGRPCError(tt.err, tt.context)

			if tt.expectedNil {
				if result != nil {
					t.Errorf("Expected nil error, got %v", result)
				}
				return
			}

			if result == nil {
				t.Error("Expected non-nil error, got nil")
				return
			}

			togomqErr, ok := result.(*TogoMQError)
			if !ok {
				t.Errorf("Expected TogoMQError, got %T", result)
				return
			}

			if togomqErr.Code != tt.expectedCode {
				t.Errorf("Expected code '%s', got '%s'", tt.expectedCode, togomqErr.Code)
			}
		})
	}
}
