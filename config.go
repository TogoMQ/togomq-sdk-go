package togomq

import (
	"fmt"
	"strings"
	"time"
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
	// UseTLS enables TLS for the connection (default: true)
	UseTLS bool
	// MaxMessageSize is the maximum message size in bytes for both send and receive (default: 50MB)
	MaxMessageSize int
	// InitialWindowSize is the initial window size for flow control (default: 128MB)
	InitialWindowSize int32
	// InitialConnWindowSize is the initial connection window size (default: 128MB)
	InitialConnWindowSize int32
	// WriteBufferSize is the write buffer size in bytes (default: 2MB)
	WriteBufferSize int
	// ReadBufferSize is the read buffer size in bytes (default: 2MB)
	ReadBufferSize int
	// KeepaliveTime is the duration after which a keepalive ping is sent (default: 60s)
	KeepaliveTime time.Duration
	// KeepaliveTimeout is the duration to wait for keepalive ping response (default: 20s)
	KeepaliveTimeout time.Duration
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	defaultMaxMessageSize := 52428800 // 50MB
	return &Config{
		Host:                  "q.togomq.io",
		Port:                  5123,
		LogLevel:              "info",
		Token:                 "",
		UseTLS:                true,
		MaxMessageSize:        defaultMaxMessageSize,
		InitialWindowSize:     128 * 1024 * 1024, // 128MB
		InitialConnWindowSize: 128 * 1024 * 1024, // 128MB
		WriteBufferSize:       2 * 1024 * 1024,   // 2MB
		ReadBufferSize:        2 * 1024 * 1024,   // 2MB
		KeepaliveTime:         60 * time.Second,  // 60s
		KeepaliveTimeout:      20 * time.Second,  // 20s
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
	if c.KeepaliveTime <= 0 {
		return fmt.Errorf("keepalive time must be greater than 0")
	}
	if c.KeepaliveTimeout <= 0 {
		return fmt.Errorf("keepalive timeout must be greater than 0")
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

// WithUseTLS sets whether to use TLS for the connection
func WithUseTLS(useTLS bool) ConfigOption {
	return func(c *Config) {
		c.UseTLS = useTLS
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

// WithKeepaliveTime sets the keepalive time duration
func WithKeepaliveTime(duration time.Duration) ConfigOption {
	return func(c *Config) {
		c.KeepaliveTime = duration
	}
}

// WithKeepaliveTimeout sets the keepalive timeout duration
func WithKeepaliveTimeout(duration time.Duration) ConfigOption {
	return func(c *Config) {
		c.KeepaliveTimeout = duration
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
