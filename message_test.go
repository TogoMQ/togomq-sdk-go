package togomq

import (
	"testing"
)

func TestNewMessage(t *testing.T) {
	topic := "test-topic"
	body := []byte("test body")

	msg := NewMessage(topic, body)

	if msg.Topic != topic {
		t.Errorf("Expected topic '%s', got '%s'", topic, msg.Topic)
	}
	if string(msg.Body) != string(body) {
		t.Errorf("Expected body '%s', got '%s'", string(body), string(msg.Body))
	}
	if msg.Variables == nil {
		t.Error("Expected variables to be initialized")
	}
}

func TestMessageBuilders(t *testing.T) {
	msg := NewMessage("test", []byte("body")).
		WithVariables(map[string]string{"key1": "value1", "key2": "value2"}).
		WithPostpone(100).
		WithRetention(3600)

	if msg.Topic != "test" {
		t.Errorf("Expected topic 'test', got '%s'", msg.Topic)
	}
	if len(msg.Variables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(msg.Variables))
	}
	if msg.Variables["key1"] != "value1" {
		t.Errorf("Expected variable key1='value1', got '%s'", msg.Variables["key1"])
	}
	if msg.Postpone != 100 {
		t.Errorf("Expected postpone 100, got %d", msg.Postpone)
	}
	if msg.Retention != 3600 {
		t.Errorf("Expected retention 3600, got %d", msg.Retention)
	}
}

func TestNewSubscribeOptions(t *testing.T) {
	topic := "test-topic"
	opts := NewSubscribeOptions(topic)

	if opts.Topic != topic {
		t.Errorf("Expected topic '%s', got '%s'", topic, opts.Topic)
	}
	if opts.Batch != 1 {
		t.Errorf("Expected batch 1, got %d", opts.Batch)
	}
	if opts.SpeedPerSec != 0 {
		t.Errorf("Expected speed 0, got %d", opts.SpeedPerSec)
	}
}

func TestSubscribeOptionsBuilders(t *testing.T) {
	opts := NewSubscribeOptions("test").
		WithBatch(10).
		WithSpeedPerSec(100)

	if opts.Topic != "test" {
		t.Errorf("Expected topic 'test', got '%s'", opts.Topic)
	}
	if opts.Batch != 10 {
		t.Errorf("Expected batch 10, got %d", opts.Batch)
	}
	if opts.SpeedPerSec != 100 {
		t.Errorf("Expected speed 100, got %d", opts.SpeedPerSec)
	}
}

func TestMessageToPubRequest(t *testing.T) {
	msg := NewMessage("test-topic", []byte("test body")).
		WithVariables(map[string]string{"key": "value"}).
		WithPostpone(50).
		WithRetention(1800)

	req := msg.toPubRequest()

	if req.Topic != msg.Topic {
		t.Errorf("Expected topic '%s', got '%s'", msg.Topic, req.Topic)
	}
	if string(req.Body) != string(msg.Body) {
		t.Errorf("Expected body '%s', got '%s'", string(msg.Body), string(req.Body))
	}
	if req.Postpone != msg.Postpone {
		t.Errorf("Expected postpone %d, got %d", msg.Postpone, req.Postpone)
	}
	if req.Retention != msg.Retention {
		t.Errorf("Expected retention %d, got %d", msg.Retention, req.Retention)
	}
	if len(req.Variables) != len(msg.Variables) {
		t.Errorf("Expected %d variables, got %d", len(msg.Variables), len(req.Variables))
	}
}

func TestSubscribeOptionsToSubRequest(t *testing.T) {
	opts := NewSubscribeOptions("test-topic").
		WithBatch(5).
		WithSpeedPerSec(50)

	req := opts.toSubRequest()

	if req.Topic != opts.Topic {
		t.Errorf("Expected topic '%s', got '%s'", opts.Topic, req.Topic)
	}
	if req.Batch != opts.Batch {
		t.Errorf("Expected batch %d, got %d", opts.Batch, req.Batch)
	}
	if req.SpeedPerSec != opts.SpeedPerSec {
		t.Errorf("Expected speed %d, got %d", opts.SpeedPerSec, req.SpeedPerSec)
	}
}
