# TogoMQ SDK for Go

[![CI](https://github.com/TogoMQ/togomq-sdk-go/actions/workflows/ci.yml/badge.svg)](https://github.com/TogoMQ/togomq-sdk-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/TogoMQ/togomq-sdk-go.svg)](https://pkg.go.dev/github.com/TogoMQ/togomq-sdk-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The official Go SDK for [TogoMQ](https://togomq.io) - a modern, high-performance message queue service. This SDK provides a simple and intuitive API for publishing and subscribing to messages using gRPC streaming.

## Features

- 🚀 **High Performance**: Built on gRPC for efficient communication
- 📡 **Streaming Support**: Native support for streaming pub/sub operations
- 🔒 **Secure**: TLS encryption and token-based authentication
- 🎯 **Simple API**: Easy-to-use client with fluent configuration
- 📝 **Comprehensive Logging**: Configurable log levels for debugging
- ⚡ **Concurrent**: Safe for concurrent use with goroutines
- ✅ **Well Tested**: Comprehensive test coverage

## Requirements

- Go 1.24 or higher
- Access to a TogoMQ server
- Valid TogoMQ authentication token

## Installation

Install the SDK using `go get`:

```bash
go get github.com/TogoMQ/togomq-sdk-go
```

## Configuration

The SDK supports flexible configuration with sensible defaults:

### Default Configuration

```go
import "github.com/TogoMQ/togomq-sdk-go"

// Create client with defaults (only token is required)
config := togomq.NewConfig(
    togomq.WithToken("your-token-here"),
)

client, err := togomq.NewClient(config)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

### Configuration Options

| Option | Default | Description |
|--------|---------|-------------|
| `Host` | `q.togomq.io` | TogoMQ server hostname |
| `Port` | `5123` | TogoMQ server port |
| `LogLevel` | `info` | Logging level (debug, info, warn, error, none) |
| `Token` | *(required)* | Authentication token |

### Custom Configuration

```go
config := togomq.NewConfig(
    togomq.WithHost("custom.togomq.io"),
    togomq.WithPort(9000),
    togomq.WithLogLevel("debug"),
    togomq.WithToken("your-token-here"),
)
```

## Usage

### Publishing Messages

**Note:** Topic name is required for all published messages. Each message must specify a topic.

#### Publishing a Batch of Messages

```go
package main

import (
    "context"
    "log"
    
    "github.com/TogoMQ/togomq-sdk-go"
)

func main() {
    // Create client
    config := togomq.NewConfig(togomq.WithToken("your-token"))
    client, err := togomq.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // Create messages - topic is required for each message
    messages := []*togomq.Message{
        togomq.NewMessage("orders", []byte("order-1")),
        togomq.NewMessage("orders", []byte("order-2")).
            WithVariables(map[string]string{
                "priority": "high",
                "customer": "12345",
            }),
        togomq.NewMessage("orders", []byte("order-3")).
            WithPostpone(60).      // Delay 60 seconds
            WithRetention(3600),   // Keep for 1 hour
    }
    
    // Publish
    resp, err := client.PubBatch(context.Background(), messages)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Published %d messages\n", resp.MessagesReceived)
}
```

#### Publishing via Channel (Streaming)

```go
func streamPublish(client *togomq.Client) {
    ctx := context.Background()
    
    // Create a channel for messages
    msgChan := make(chan *togomq.Message, 10)
    
    // Start publishing in background
    go func() {
        resp, err := client.Pub(ctx, msgChan)
        if err != nil {
            log.Printf("Publish error: %v", err)
            return
        }
        log.Printf("Published %d messages\n", resp.MessagesReceived)
    }()
    
    // Send messages
    for i := 0; i < 100; i++ {
        msg := togomq.NewMessage("events", []byte(fmt.Sprintf("event-%d", i)))
        msgChan <- msg
    }
    
    // Close channel to signal end of stream
    close(msgChan)
    
    // Wait for completion
    time.Sleep(2 * time.Second)
}
```

### Subscribing to Messages

**Note:** Topic is required for subscriptions. Use wildcards like `"orders.*"` for pattern matching, or `"*"` to receive messages from all topics.

#### Basic Subscription

```go
package main

import (
    "context"
    "log"
    
    "github.com/TogoMQ/togomq-sdk-go"
)

func main() {
    // Create client
    config := togomq.NewConfig(togomq.WithToken("your-token"))
    client, err := togomq.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // Subscribe to specific topic
    // Topic is required - use "*" to subscribe to all topics
    opts := togomq.NewSubscribeOptions("orders")
    msgChan, errChan, err := client.Sub(context.Background(), opts)
    if err != nil {
        log.Fatal(err)
    }
    
    // Receive messages
    for {
        select {
        case msg, ok := <-msgChan:
            if !ok {
                log.Println("Subscription ended")
                return
            }
            log.Printf("Received message from %s: %s\n", msg.Topic, string(msg.Body))
            log.Printf("Message UUID: %s\n", msg.UUID)
            
            // Access variables
            if priority, ok := msg.Variables["priority"]; ok {
                log.Printf("Priority: %s\n", priority)
            }
            
        case err := <-errChan:
            log.Printf("Subscription error: %v\n", err)
            return
        }
    }
}
```

#### Advanced Subscription with Options

```go
// Subscribe with batch size and rate limiting
// Default values: Batch = 0 (default 1000 if not set), SpeedPerSec = 0 (unlimited)
opts := togomq.NewSubscribeOptions("orders.*").  // Wildcard topic
    WithBatch(10).                                // Receive up to 10 messages at once
    WithSpeedPerSec(100)                          // Limit to 100 messages per second

msgChan, errChan, err := client.Sub(context.Background(), opts)
if err != nil {
    log.Fatal(err)
}
```

**Subscription Options:**
- **Batch**: Maximum number of messages to receive at once (default: 0 = default 1000 if not set)
- **SpeedPerSec**: Rate limit for message delivery per second (default: 0 = unlimited)

#### Subscribe to All Topics (Wildcard)

```go
// Subscribe to all topics using "*" wildcard
opts := togomq.NewSubscribeOptions("*") // "*" = all topics
msgChan, errChan, err := client.Sub(ctx, opts)
```

#### Subscribe with Pattern Wildcards

```go
// Subscribe to all orders topics (orders.new, orders.updated, etc.)
opts := togomq.NewSubscribeOptions("orders.*")
msgChan, errChan, err := client.Sub(ctx, opts)

// Subscribe to all topics
opts := togomq.NewSubscribeOptions("*")
msgChan, errChan, err := client.Sub(ctx, opts)
```

#### Subscription with Context Cancellation

```go
// Create cancellable context
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

opts := togomq.NewSubscribeOptions("events")
msgChan, errChan, err := client.Sub(ctx, opts)
if err != nil {
    log.Fatal(err)
}

// Process messages until context is cancelled
for {
    select {
    case msg := <-msgChan:
        if msg == nil {
            return
        }
        processMessage(msg)
    case err := <-errChan:
        log.Printf("Error: %v\n", err)
        return
    case <-ctx.Done():
        log.Println("Context cancelled")
        return
    }
}
```

## Message Structure

### Publishing Message

**Important:** Topic is required when publishing messages.

```go
type Message struct {
    Topic     string            // Message topic (required)
    Body      []byte            // Message payload
    Variables map[string]string // Custom key-value metadata
    Postpone  int64             // Delay in seconds before message is available
    Retention int64             // How long to keep message (seconds)
}
```

### Received Message

```go
type Message struct {
    Topic     string            // Message topic
    UUID      string            // Unique message identifier
    Body      []byte            // Message payload
    Variables map[string]string // Custom key-value metadata
}
```

## Error Handling

The SDK provides detailed error information:

```go
resp, err := client.PubBatch(ctx, messages)
if err != nil {
    // Check error type
    if togomqErr, ok := err.(*togomq.TogoMQError); ok {
        switch togomqErr.Code {
        case togomq.ErrCodeAuth:
            log.Println("Authentication failed")
        case togomq.ErrCodeConnection:
            log.Println("Connection error")
        case togomq.ErrCodeValidation:
            log.Println("Validation error")
        default:
            log.Printf("Error: %v\n", togomqErr)
        }
    }
}
```

### Error Codes

- `ErrCodeConnection` - Connection or network errors
- `ErrCodeAuth` - Authentication failures
- `ErrCodeValidation` - Invalid input or configuration
- `ErrCodePublish` - Publishing errors
- `ErrCodeSubscribe` - Subscription errors
- `ErrCodeStream` - General streaming errors
- `ErrCodeConfiguration` - Configuration errors

## Logging

Control logging verbosity with the `LogLevel` configuration:

```go
config := togomq.NewConfig(
    togomq.WithToken("your-token"),
    togomq.WithLogLevel("debug"), // debug, info, warn, error, none
)
```

Log levels:
- `debug` - All logs including debug information
- `info` - Informational messages and above
- `warn` - Warnings and errors only
- `error` - Error messages only
- `none` - Disable logging

## Best Practices

1. **Reuse Clients**: Create one client per application and reuse it across goroutines
2. **Handle Errors**: Always check and handle errors appropriately
3. **Close Connections**: Always defer `client.Close()` after creating a client
4. **Use Context**: Leverage context for timeouts and cancellation
5. **Batch Messages**: Use `PubBatch` for better performance when publishing multiple messages
6. **Monitor Channels**: Always monitor both message and error channels in subscriptions

## Examples

Check out the `examples/` directory for complete working examples:

- `examples/publish/` - Publishing examples
- `examples/subscribe/` - Subscription examples
- `examples/advanced/` - Advanced usage patterns

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- Documentation: [https://togomq.io/docs](https://togomq.io/docs)
- Issues: [https://github.com/TogoMQ/togomq-sdk-go/issues](https://github.com/TogoMQ/togomq-sdk-go/issues)
- TogoMQ Website: [https://togomq.io](https://togomq.io)

## Related Projects

- [togomq-grpc-go](https://github.com/TogoMQ/togomq-grpc-go) - Auto-generated gRPC protobuf definitions
