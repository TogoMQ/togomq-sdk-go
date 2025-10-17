package togomq

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error codes
const (
	ErrCodeConnection    = "CONNECTION_ERROR"
	ErrCodeAuth          = "AUTH_ERROR"
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodePublish       = "PUBLISH_ERROR"
	ErrCodeSubscribe     = "SUBSCRIBE_ERROR"
	ErrCodeStream        = "STREAM_ERROR"
	ErrCodeConfiguration = "CONFIG_ERROR"
)

// TogoMQError represents an error from the TogoMQ SDK
type TogoMQError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *TogoMQError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *TogoMQError) Unwrap() error {
	return e.Err
}

// NewError creates a new TogoMQError
func NewError(code, message string, err error) *TogoMQError {
	return &TogoMQError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WrapGRPCError wraps a gRPC error with context
func WrapGRPCError(err error, context string) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return NewError(ErrCodeStream, context, err)
	}

	var code string
	switch st.Code() {
	case codes.Unauthenticated:
		code = ErrCodeAuth
	case codes.InvalidArgument:
		code = ErrCodeValidation
	case codes.Unavailable:
		code = ErrCodeConnection
	default:
		code = ErrCodeStream
	}

	return NewError(code, fmt.Sprintf("%s: %s", context, st.Message()), err)
}
