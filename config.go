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
	// MaxMessageSize is the maximum message size in bytes for both send and receive (default: 50MB)
	MaxMessageSize int
	// InitialWindowSize is the initial window size for flow control (default: same as MaxMessageSize)
	InitialWindowSize int32
	// InitialConnWindowSize is the initial connection window size (default: same as MaxMessageSize)
	InitialConnWindowSize int32
	// WriteBufferSize is the write buffer size in bytes (default: 256KB)
	WriteBufferSize int
	// ReadBufferSize is the read buffer size in bytes (default: 256KB)
	ReadBufferSize int
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	defaultMaxMessageSize := 52428800 // 50MB
	return &Config{
		Host:                  "q.togomq.io",
		Port:                  5123,
		LogLevel:              "info",
		Token:                 "",
		MaxMessageSize:        defaultMaxMessageSize,
		InitialWindowSize:     int32(defaultMaxMessageSize),
		InitialConnWindowSize: int32(defaultMaxMessageSize),
		WriteBufferSize:       256 * 1024, // 256KB
		ReadBufferSize:        256 * 1024, // 256KB
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
	if c.MaxMessageSize <= 0 {
		return fmt.Errorf("max message size must be greater than 0")
	}
	if c.InitialWindowSize <= 0 {
		return fmt.Errorf("initial window size must be greater than 0")
	}
	if c.InitialConnWindowSize <= 0 {
		return fmt.Errorf("initial connection window size must be greater than 0")
	}
	if c.WriteBufferSize <= 0 {
		return fmt.Errorf("write buffer size must be greater than 0")
	}
	if c.ReadBufferSize <= 0 {
		return fmt.Errorf("read buffer size must be greater than 0")
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

// WithMaxMessageSize sets the maximum message size in bytes
func WithMaxMessageSize(size int) ConfigOption {
	return func(c *Config) {
		c.MaxMessageSize = size
		// Also update window sizes to match
		c.InitialWindowSize = int32(size)
		c.InitialConnWindowSize = int32(size)
	}
}

// WithInitialWindowSize sets the initial window size for flow control
func WithInitialWindowSize(size int32) ConfigOption {
	return func(c *Config) {
		c.InitialWindowSize = size
	}
}

// WithInitialConnWindowSize sets the initial connection window size
func WithInitialConnWindowSize(size int32) ConfigOption {
	return func(c *Config) {
		c.InitialConnWindowSize = size
	}
}

// WithWriteBufferSize sets the write buffer size in bytes
func WithWriteBufferSize(size int) ConfigOption {
	return func(c *Config) {
		c.WriteBufferSize = size
	}
}

// WithReadBufferSize sets the read buffer size in bytes
func WithReadBufferSize(size int) ConfigOption {
	return func(c *Config) {
		c.ReadBufferSize = size
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
