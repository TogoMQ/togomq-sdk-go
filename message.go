package togomq

import (
	mqv1 "github.com/TogoMQ/togomq-grpc-go/mq/v1"
)

// Message represents a message to be published or received
type Message struct {
	// Topic is the message topic
	Topic string
	// Body is the message payload
	Body []byte
	// Variables are custom key-value pairs associated with the message
	Variables map[string]string
	// Postpone is the delay in seconds before the message becomes available (for publishing)
	Postpone int64
	// Retention is how long the message should be kept in seconds (for publishing)
	Retention int64
	// UUID is the unique identifier of the message (for received messages)
	UUID string
}

// NewMessage creates a new message with the given topic and body
func NewMessage(topic string, body []byte) *Message {
	return &Message{
		Topic:     topic,
		Body:      body,
		Variables: make(map[string]string),
	}
}

// WithVariables adds variables to the message
func (m *Message) WithVariables(vars map[string]string) *Message {
	m.Variables = vars
	return m
}

// WithPostpone sets the postpone delay
func (m *Message) WithPostpone(postpone int64) *Message {
	m.Postpone = postpone
	return m
}

// WithRetention sets the retention period
func (m *Message) WithRetention(retention int64) *Message {
	m.Retention = retention
	return m
}

// toPubRequest converts a Message to a gRPC PubMessageRequest
func (m *Message) toPubRequest() *mqv1.PubMessageRequest {
	return &mqv1.PubMessageRequest{
		Topic:     m.Topic,
		Body:      m.Body,
		Variables: m.Variables,
		Postpone:  m.Postpone,
		Retention: m.Retention,
	}
}

// fromSubResponse converts a gRPC SubMessageResponse to a Message
func fromSubResponse(resp *mqv1.SubMessageResponse) *Message {
	return &Message{
		Topic:     resp.Topic,
		UUID:      resp.Uuid,
		Body:      resp.Body,
		Variables: resp.Variables,
	}
}

// SubscribeOptions represents options for subscribing to messages.
type SubscribeOptions struct {
	// Topic is the topic to subscribe to (required, supports wildcards like "orders.*" or "*" for all topics)
	Topic string
	// Batch is the maximum number of messages to receive at once (0 = default 1000 if not set)
	Batch int64
	// SpeedPerSec limits the rate of message delivery per second (0 = unlimited)
	SpeedPerSec int64
}

// NewSubscribeOptions creates default subscribe options
func NewSubscribeOptions(topic string) *SubscribeOptions {
	return &SubscribeOptions{
		Topic:       topic,
		Batch:       0, // default 1000 if not set
		SpeedPerSec: 0, // unlimited
	}
}

// WithBatch sets the batch size
func (s *SubscribeOptions) WithBatch(batch int64) *SubscribeOptions {
	s.Batch = batch
	return s
}

// WithSpeedPerSec sets the speed limit
func (s *SubscribeOptions) WithSpeedPerSec(speed int64) *SubscribeOptions {
	s.SpeedPerSec = speed
	return s
}

// toSubRequest converts SubscribeOptions to a gRPC SubMessageRequest
func (s *SubscribeOptions) toSubRequest() *mqv1.SubMessageRequest {
	return &mqv1.SubMessageRequest{
		Topic:       s.Topic,
		Batch:       s.Batch,
		SpeedPerSec: s.SpeedPerSec,
	}
}

// PubResponse contains the result of a publish operation
type PubResponse struct {
	// MessagesReceived is the number of messages successfully received by the server
	MessagesReceived int64
}
