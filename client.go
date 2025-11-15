package togomq

import (
	"context"
	"io"

	mqv1 "github.com/TogoMQ/togomq-grpc-go/mq/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

// Client is the TogoMQ client
type Client struct {
	config *Config
	conn   *grpc.ClientConn
	client mqv1.MqServiceClient
	logger *Logger
}

// NewClient creates a new TogoMQ client
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, NewError(ErrCodeConfiguration, "config cannot be nil", nil)
	}

	if err := config.Validate(); err != nil {
		return nil, NewError(ErrCodeValidation, "invalid configuration", err)
	}

	logger := NewLogger(ParseLogLevel(config.LogLevel))
	logger.Info("Creating TogoMQ client for %s", config.Address())

	// Create TLS credentials
	creds := credentials.NewTLS(nil)

	// Create gRPC connection with performance and large message support settings
	conn, err := grpc.NewClient(
		config.Address(),
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(config.MaxMessageSize),
			grpc.MaxCallSendMsgSize(config.MaxMessageSize),
		),
		grpc.WithInitialWindowSize(config.InitialWindowSize),
		grpc.WithInitialConnWindowSize(config.InitialConnWindowSize),
		grpc.WithWriteBufferSize(config.WriteBufferSize),
		grpc.WithReadBufferSize(config.ReadBufferSize),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                config.KeepaliveTime,
			Timeout:             config.KeepaliveTimeout,
			PermitWithoutStream: false,
		}),
	)
	if err != nil {
		logger.Error("Failed to connect to TogoMQ: %v", err)
		return nil, NewError(ErrCodeConnection, "failed to create gRPC connection", err)
	}

	client := mqv1.NewMqServiceClient(conn)

	logger.Info("TogoMQ client created successfully")

	return &Client{
		config: config,
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	c.logger.Info("Closing TogoMQ client")
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// contextWithAuth adds authentication metadata to the context
func (c *Client) contextWithAuth(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{
		"authorization": c.config.Token,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

// Pub publishes messages to TogoMQ using a streaming approach
// Messages are sent through the provided channel and the function returns when the channel is closed
func (c *Client) Pub(ctx context.Context, messages <-chan *Message) (*PubResponse, error) {
	c.logger.Debug("Starting Pub operation")

	// Add authentication to context
	ctx = c.contextWithAuth(ctx)

	// Create the stream
	stream, err := c.client.PubMessage(ctx)
	if err != nil {
		c.logger.Error("Failed to create pub stream: %v", err)
		return nil, WrapGRPCError(err, "failed to create publish stream")
	}

	// Send messages
	messageCount := 0
	for msg := range messages {
		// Validate that topic is specified
		if msg.Topic == "" {
			c.logger.Error("Message topic is required")
			return nil, NewError(ErrCodeValidation, "message topic is required", nil)
		}

		c.logger.Debug("Publishing message to topic: %s", msg.Topic)

		if err := stream.Send(msg.toPubRequest()); err != nil {
			c.logger.Error("Failed to send message: %v", err)
			return nil, WrapGRPCError(err, "failed to send message")
		}
		messageCount++
	}

	c.logger.Info("Sent %d messages, waiting for response", messageCount)

	// Close and receive response
	resp, err := stream.CloseAndRecv()
	if err != nil {
		c.logger.Error("Failed to receive pub response: %v", err)
		return nil, WrapGRPCError(err, "failed to receive publish response")
	}

	c.logger.Info("Publish completed: %d messages received by server", resp.MessagesReceived)

	return &PubResponse{
		MessagesReceived: resp.MessagesReceived,
	}, nil
}

// PubBatch publishes a batch of messages
func (c *Client) PubBatch(ctx context.Context, messages []*Message) (*PubResponse, error) {
	c.logger.Debug("Publishing batch of %d messages", len(messages))

	// Create a channel and send messages
	msgChan := make(chan *Message, len(messages))
	for _, msg := range messages {
		msgChan <- msg
	}
	close(msgChan)

	return c.Pub(ctx, msgChan)
}

// Sub subscribes to messages from TogoMQ.
// Topic is required (can use wildcards like "orders.*" or "*" for all topics).
// Returns channels for messages and errors, and an error if the subscription fails to start.
func (c *Client) Sub(ctx context.Context, opts *SubscribeOptions) (<-chan *Message, <-chan error, error) {
	// Validate that topic is specified
	if opts.Topic == "" {
		return nil, nil, NewError(ErrCodeValidation, "topic is required for subscription", nil)
	}

	c.logger.Debug("Starting Sub operation for topic: %s", opts.Topic)
	// Add authentication to context
	ctx = c.contextWithAuth(ctx)

	// Create the stream
	stream, err := c.client.SubMessage(ctx, opts.toSubRequest())
	if err != nil {
		c.logger.Error("Failed to create sub stream: %v", err)
		return nil, nil, WrapGRPCError(err, "failed to create subscribe stream")
	}

	// Create channels for messages and errors
	messageChan := make(chan *Message)
	errorChan := make(chan error, 1)

	// Start goroutine to receive messages
	go func() {
		defer close(messageChan)
		defer close(errorChan)

		messageCount := 0
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				c.logger.Info("Subscribe stream ended, received %d messages", messageCount)
				return
			}
			if err != nil {
				c.logger.Error("Failed to receive message: %v", err)
				errorChan <- WrapGRPCError(err, "failed to receive message")
				return
			}

			c.logger.Debug("Received message from topic: %s, UUID: %s", resp.Topic, resp.Uuid)
			messageCount++

			msg := fromSubResponse(resp)

			select {
			case messageChan <- msg:
				// Message sent successfully
			case <-ctx.Done():
				c.logger.Info("Context cancelled, stopping subscription")
				return
			}
		}
	}()

	if opts.Topic == "" {
		c.logger.Info("Subscribe stream started for all topics (wildcard)")
	} else {
		c.logger.Info("Subscribe stream started for topic: %s", opts.Topic)
	}

	return messageChan, errorChan, nil
}

// CountMessages counts the number of messages in a topic.
// Topic can use wildcards (e.g., "orders.*" or "*" for all topics).
// Returns the total count of messages matching the topic pattern.
func (c *Client) CountMessages(ctx context.Context, topic string) (int64, error) {
	// Validate that topic is specified
	if topic == "" {
		return 0, NewError(ErrCodeValidation, "topic is required for counting messages", nil)
	}

	c.logger.Debug("Counting messages for topic: %s", topic)

	// Add authentication to context
	ctx = c.contextWithAuth(ctx)

	// Create the request
	req := &mqv1.CountMessagesRequest{
		Topic: topic,
	}

	// Call the gRPC method
	resp, err := c.client.CountMessages(ctx, req)
	if err != nil {
		c.logger.Error("Failed to count messages: %v", err)
		return 0, WrapGRPCError(err, "failed to count messages")
	}

	c.logger.Info("Counted %d messages for topic: %s", resp.MessagesCount, topic)

	return resp.MessagesCount, nil
}
