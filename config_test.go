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
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: &Config{
				Host:     "test.example.com",
				Port:     5123,
				LogLevel: "info",
				Token:    "mytoken",
			},
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
