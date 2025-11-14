package togomq

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Host != "q.togomq.io" {
		t.Errorf("Expected default host to be 'q.togomq.io', got '%s'", cfg.Host)
	}
	if cfg.Port != 5123 {
		t.Errorf("Expected default port to be 5123, got %d", cfg.Port)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("Expected default log level to be 'info', got '%s'", cfg.LogLevel)
	}
	if cfg.Token != "" {
		t.Errorf("Expected default token to be empty, got '%s'", cfg.Token)
	}
	// Check gRPC settings
	expectedMaxMsgSize := 52428800 // 50MB
	if cfg.MaxMessageSize != expectedMaxMsgSize {
		t.Errorf("Expected default max message size to be %d, got %d", expectedMaxMsgSize, cfg.MaxMessageSize)
	}
	if cfg.InitialWindowSize != int32(expectedMaxMsgSize) {
		t.Errorf("Expected default initial window size to be %d, got %d", expectedMaxMsgSize, cfg.InitialWindowSize)
	}
	if cfg.InitialConnWindowSize != int32(expectedMaxMsgSize) {
		t.Errorf("Expected default initial conn window size to be %d, got %d", expectedMaxMsgSize, cfg.InitialConnWindowSize)
	}
	if cfg.WriteBufferSize != 256*1024 {
		t.Errorf("Expected default write buffer size to be 262144, got %d", cfg.WriteBufferSize)
	}
	if cfg.ReadBufferSize != 256*1024 {
		t.Errorf("Expected default read buffer size to be 262144, got %d", cfg.ReadBufferSize)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid config",
			config:      NewConfig(WithToken("mytoken")),
			expectError: false,
		},
		{
			name: "empty host",
			config: &Config{
				Host:  "",
				Port:  5123,
				Token: "mytoken",
			},
			expectError: true,
			errorMsg:    "host cannot be empty",
		},
		{
			name: "invalid port - zero",
			config: &Config{
				Host:  "test.example.com",
				Port:  0,
				Token: "mytoken",
			},
			expectError: true,
			errorMsg:    "port must be between 1 and 65535",
		},
		{
			name: "invalid port - too high",
			config: &Config{
				Host:  "test.example.com",
				Port:  70000,
				Token: "mytoken",
			},
			expectError: true,
			errorMsg:    "port must be between 1 and 65535",
		},
		{
			name: "empty token",
			config: &Config{
				Host:  "test.example.com",
				Port:  5123,
				Token: "",
			},
			expectError: true,
			errorMsg:    "token is required",
		},
		{
			name: "invalid max message size",
			config: &Config{
				Host:           "test.example.com",
				Port:           5123,
				Token:          "mytoken",
				MaxMessageSize: 0,
			},
			expectError: true,
			errorMsg:    "max message size must be greater than 0",
		},
		{
			name: "invalid initial window size",
			config: &Config{
				Host:              "test.example.com",
				Port:              5123,
				Token:             "mytoken",
				MaxMessageSize:    1024,
				InitialWindowSize: 0,
			},
			expectError: true,
			errorMsg:    "initial window size must be greater than 0",
		},
		{
			name: "invalid initial conn window size",
			config: &Config{
				Host:                  "test.example.com",
				Port:                  5123,
				Token:                 "mytoken",
				MaxMessageSize:        1024,
				InitialWindowSize:     1024,
				InitialConnWindowSize: 0,
			},
			expectError: true,
			errorMsg:    "initial connection window size must be greater than 0",
		},
		{
			name: "invalid write buffer size",
			config: &Config{
				Host:                  "test.example.com",
				Port:                  5123,
				Token:                 "mytoken",
				MaxMessageSize:        1024,
				InitialWindowSize:     1024,
				InitialConnWindowSize: 1024,
				WriteBufferSize:       0,
			},
			expectError: true,
			errorMsg:    "write buffer size must be greater than 0",
		},
		{
			name: "invalid read buffer size",
			config: &Config{
				Host:                  "test.example.com",
				Port:                  5123,
				Token:                 "mytoken",
				MaxMessageSize:        1024,
				InitialWindowSize:     1024,
				InitialConnWindowSize: 1024,
				WriteBufferSize:       1024,
				ReadBufferSize:        0,
			},
			expectError: true,
			errorMsg:    "read buffer size must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error '%s', got nil", tt.errorMsg)
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got '%v'", err)
				}
			}
		})
	}
}

func TestConfigAddress(t *testing.T) {
	cfg := &Config{
		Host: "example.com",
		Port: 8080,
	}

	expected := "example.com:8080"
	if cfg.Address() != expected {
		t.Errorf("Expected address '%s', got '%s'", expected, cfg.Address())
	}
}

func TestConfigOptions(t *testing.T) {
	cfg := NewConfig(
		WithHost("custom.example.com"),
		WithPort(9000),
		WithLogLevel("debug"),
		WithToken("custom-token"),
		WithMaxMessageSize(10485760), // 10MB
		WithWriteBufferSize(512*1024),
		WithReadBufferSize(512*1024),
	)

	if cfg.Host != "custom.example.com" {
		t.Errorf("Expected host 'custom.example.com', got '%s'", cfg.Host)
	}
	if cfg.Port != 9000 {
		t.Errorf("Expected port 9000, got %d", cfg.Port)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", cfg.LogLevel)
	}
	if cfg.Token != "custom-token" {
		t.Errorf("Expected token 'custom-token', got '%s'", cfg.Token)
	}
	if cfg.MaxMessageSize != 10485760 {
		t.Errorf("Expected max message size 10485760, got %d", cfg.MaxMessageSize)
	}
	// WithMaxMessageSize should also set window sizes
	if cfg.InitialWindowSize != 10485760 {
		t.Errorf("Expected initial window size 10485760, got %d", cfg.InitialWindowSize)
	}
	if cfg.InitialConnWindowSize != 10485760 {
		t.Errorf("Expected initial conn window size 10485760, got %d", cfg.InitialConnWindowSize)
	}
	if cfg.WriteBufferSize != 512*1024 {
		t.Errorf("Expected write buffer size 524288, got %d", cfg.WriteBufferSize)
	}
	if cfg.ReadBufferSize != 512*1024 {
		t.Errorf("Expected read buffer size 524288, got %d", cfg.ReadBufferSize)
	}
}

func TestWindowSizeOptions(t *testing.T) {
	cfg := NewConfig(
		WithToken("test-token"),
		WithInitialWindowSize(1024),
		WithInitialConnWindowSize(2048),
	)

	if cfg.InitialWindowSize != 1024 {
		t.Errorf("Expected initial window size 1024, got %d", cfg.InitialWindowSize)
	}
	if cfg.InitialConnWindowSize != 2048 {
		t.Errorf("Expected initial conn window size 2048, got %d", cfg.InitialConnWindowSize)
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"debug", LogLevelDebug},
		{"DEBUG", LogLevelDebug},
		{"info", LogLevelInfo},
		{"INFO", LogLevelInfo},
		{"warn", LogLevelWarn},
		{"warning", LogLevelWarn},
		{"error", LogLevelError},
		{"ERROR", LogLevelError},
		{"none", LogLevelNone},
		{"invalid", LogLevelInfo}, // defaults to info
		{"", LogLevelInfo},        // defaults to info
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseLogLevel(tt.input)
			if result != tt.expected {
				t.Errorf("ParseLogLevel(%s) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}
