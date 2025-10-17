package main

import (
	"context"
	"log"
	"time"

	"github.com/TogoMQ/togomq-sdk-go"
)

func main() {
	// Example: Publishing messages to TogoMQ
	publishExample()

	// Example: Subscribing to messages from TogoMQ
	// subscribeExample()
}

func publishExample() {
	// Create client configuration
	config := togomq.NewConfig(
		togomq.WithToken("your-token-here"), // Replace with your actual token
		togomq.WithLogLevel("info"),
	)

	// Create client
	client, err := togomq.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Create messages
	messages := []*togomq.Message{
		togomq.NewMessage("orders", []byte("order-1-data")),
		togomq.NewMessage("orders", []byte("order-2-data")).
			WithVariables(map[string]string{
				"priority": "high",
				"customer": "12345",
			}),
		togomq.NewMessage("orders", []byte("order-3-data")).
			WithPostpone(60).    // Delay 60 seconds
			WithRetention(3600), // Keep for 1 hour
	}

	// Publish messages
	ctx := context.Background()
	resp, err := client.PubBatch(ctx, messages)
	if err != nil {
		log.Fatalf("Failed to publish messages: %v", err)
	}

	log.Printf("Successfully published %d messages\n", resp.MessagesReceived)
}

func subscribeExample() {
	// Create configuration with token and log level
	config := togomq.NewConfig(
		togomq.WithToken("token"),
		togomq.WithLogLevel("info"),
	)

	// Create client
	client, err := togomq.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Subscribe to messages
	// Topic is required - use "*" to subscribe to all topics, or "orders.*" for pattern matching
	opts := togomq.NewSubscribeOptions("orders")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msgChan, errChan, err := client.Sub(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	log.Println("Listening for messages...")

	// Process messages
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				log.Println("Subscription ended")
				return
			}
			log.Printf("Received message from %s: %s\n", msg.Topic, string(msg.Body))
			if len(msg.Variables) > 0 {
				log.Printf("Variables: %+v\n", msg.Variables)
			}

		case err := <-errChan:
			log.Printf("Subscription error: %v\n", err)
			return

		case <-ctx.Done():
			log.Println("Context cancelled, ending subscription")
			return
		}
	}
}
