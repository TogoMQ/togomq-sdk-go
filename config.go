package togomq

import (
	"fmt"
	"strings"
)

// Config holds the configuration for the TogoMQ client
type Config struct {
	// Host is the TogoMQ server hostname
	Host string
	// Port is the TogoMQ server port
	Port int
	// LogLevel defines the logging verbosity
	LogLevel string
	// Token is the authentication token (required)
	Token string
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		Host:     "q.togomq.io",
		Port:     5123,
		LogLevel: "info",
		Token:    "",
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if strings.TrimSpace(c.Host) == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if strings.TrimSpace(c.Token) == "" {
		return fmt.Errorf("token is required")
	}
	return nil
}

// Address returns the full server address
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// ConfigOption is a function that modifies a Config
type ConfigOption func(*Config)

// WithHost sets the host
func WithHost(host string) ConfigOption {
	return func(c *Config) {
		c.Host = host
	}
}

// WithPort sets the port
func WithPort(port int) ConfigOption {
	return func(c *Config) {
		c.Port = port
	}
}

// WithLogLevel sets the log level
func WithLogLevel(level string) ConfigOption {
	return func(c *Config) {
		c.LogLevel = level
	}
}

// WithToken sets the token
func WithToken(token string) ConfigOption {
	return func(c *Config) {
		c.Token = token
	}
}

// NewConfig creates a new Config with optional overrides
func NewConfig(opts ...ConfigOption) *Config {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
