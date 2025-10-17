# TogoMQ SDK Examples

This directory contains complete working examples demonstrating various features of the TogoMQ SDK for Go.

## Prerequisites

Before running any examples:

1. Install the TogoMQ SDK:
   ```bash
   go get github.com/TogoMQ/togomq-sdk-go
   ```

2. Get a TogoMQ authentication token from [https://togomq.io](https://togomq.io)

3. Replace `"your-token-here"` in the example code with your actual token

## Examples Directory Structure

### 📤 `publish/` - Publishing Examples

Demonstrates different ways to publish messages to TogoMQ:

- **Batch Publishing** - Publishing multiple messages at once for better performance
- **Streaming Publishing** - Publishing messages via a channel for real-time scenarios

**Run:**
```bash
cd publish
go run main.go
```

**Features demonstrated:**
- Creating messages with topic (required)
- Adding custom variables to messages
- Using postpone for delayed message delivery
- Setting retention time for messages
- Batch vs streaming publishing patterns

---

### 📥 `subscribe/` - Subscription Examples

Shows various subscription patterns and options:

- **Basic Subscription** - Subscribe to a specific topic
- **Wildcard Subscription** - Subscribe to all topics using `"*"`
- **Pattern Matching** - Subscribe using patterns like `"orders.*"`
- **Advanced Subscription** - Using batch size and rate limiting options

**Run:**
```bash
cd subscribe
go run main.go
```

**Note:** By default, the `advancedSubscribe` example is active. Uncomment other functions in `main()` to try different subscription patterns.

**Features demonstrated:**
- Topic-specific subscriptions
- Wildcard and pattern matching
- Message batching
- Rate limiting with `WithSpeedPerSec()`
- Context cancellation and timeouts
- Accessing message UUID and variables

---

### 🚀 `advanced/` - Advanced Usage Patterns

Complex scenarios and best practices:

- **Concurrent Pub/Sub** - Multiple publishers and subscribers running simultaneously
- **Error Handling** - Comprehensive error handling with TogoMQ error codes
- **Custom Configuration** - Using custom host, port, and logging options
- **Full Featured Messages** - Messages using all available features (variables, postpone, retention)

**Run:**
```bash
cd advanced
go run main.go
```

**Note:** By default, the `customConfigExample` is active. Uncomment other functions in `main()` to try different advanced patterns.

**Features demonstrated:**
- Concurrent operations with goroutines
- Proper error handling with error codes
- Custom client configuration
- Using all message features together
- Wait groups for synchronization
- Production-ready patterns

---

## Quick Start

### Publishing Messages

```go
// Create client
config := togomq.NewConfig(togomq.WithToken("your-token"))
client, err := togomq.NewClient(config)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// Create and publish message
msg := togomq.NewMessage("my-topic", []byte("Hello TogoMQ!"))
resp, err := client.PubBatch(context.Background(), []*togomq.Message{msg})
if err != nil {
    log.Fatal(err)
}
log.Printf("Published %d messages\n", resp.MessagesReceived)
```

### Subscribing to Messages

```go
// Create client
config := togomq.NewConfig(togomq.WithToken("your-token"))
client, err := togomq.NewClient(config)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// Subscribe to topic
opts := togomq.NewSubscribeOptions("my-topic")
msgChan, errChan, err := client.Sub(context.Background(), opts)
if err != nil {
    log.Fatal(err)
}

// Receive messages
for {
    select {
    case msg := <-msgChan:
        log.Printf("Received: %s\n", string(msg.Body))
    case err := <-errChan:
        log.Printf("Error: %v\n", err)
        return
    }
}
```

## Common Patterns

### Using Context for Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

opts := togomq.NewSubscribeOptions("my-topic")
msgChan, errChan, err := client.Sub(ctx, opts)
// ... handle messages with timeout
```

### Message with Variables

```go
msg := togomq.NewMessage("orders", []byte(`{"id": 123}`)).
    WithVariables(map[string]string{
        "priority": "high",
        "customer": "CUST-456",
    })
```

### Delayed Messages

```go
// Delay message by 5 minutes
msg := togomq.NewMessage("scheduled", []byte("task data")).
    WithPostpone(300) // 300 seconds = 5 minutes
```

### Error Handling

```go
resp, err := client.PubBatch(ctx, messages)
if err != nil {
    if togomqErr, ok := err.(*togomq.TogoMQError); ok {
        switch togomqErr.Code {
        case togomq.ErrCodeAuth:
            log.Println("Authentication failed")
        case togomq.ErrCodeValidation:
            log.Println("Invalid input")
        default:
            log.Printf("Error: %v\n", togomqErr)
        }
    }
}
```

## Tips for Development

1. **Start Simple**: Begin with the `publish/` example to understand basic publishing
2. **Topic Names**: Always specify a topic when publishing; use `"*"` to subscribe to all topics
3. **Error Handling**: Always check and handle errors appropriately
4. **Close Connections**: Use `defer client.Close()` to ensure proper cleanup
5. **Context Usage**: Use context for timeouts and cancellation in production code
6. **Reuse Clients**: Create one client and reuse it across your application

## Troubleshooting

**Connection Errors:**
- Verify your token is correct
- Check that you can reach `q.togomq.io:5123`
- Try increasing log level: `WithLogLevel("debug")`

**No Messages Received:**
- Ensure publisher and subscriber use matching topics
- Check that messages aren't postponed
- Verify subscriber is running before publishing

**Authentication Errors:**
- Confirm your token is valid
- Check token hasn't expired
- Ensure token has proper permissions

## Next Steps

- Read the [main README](../README.md) for complete API reference
- Check the [AGENTS.md](../AGENTS.md) for development guidelines
- Visit [TogoMQ Documentation](https://togomq.io/docs) for server-side details

## Support

- Issues: [GitHub Issues](https://github.com/TogoMQ/togomq-sdk-go/issues)
- Documentation: [https://togomq.io/docs](https://togomq.io/docs)
