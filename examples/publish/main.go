package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TogoMQ/togomq-sdk-go"
)

func main() {
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

	// Example 1: Publishing a batch of messages
	fmt.Println("=== Example 1: Batch Publishing ===")
	publishBatch(client)

	// Example 2: Publishing via streaming channel
	fmt.Println("\n=== Example 2: Streaming Publishing ===")
	publishStream(client)
}

// publishBatch demonstrates publishing a batch of messages at once
func publishBatch(client *togomq.Client) {
	// Create messages - topic is required for each message
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

// publishStream demonstrates publishing messages via a channel (streaming)
func publishStream(client *togomq.Client) {
	ctx := context.Background()

	// Create a channel for messages
	msgChan := make(chan *togomq.Message, 10)

	// Start publishing in background
	done := make(chan bool)
	go func() {
		resp, err := client.Pub(ctx, msgChan)
		if err != nil {
			log.Printf("Publish error: %v", err)
			done <- false
			return
		}
		log.Printf("Published %d messages via streaming\n", resp.MessagesReceived)
		done <- true
	}()

	// Send messages through the channel
	for i := 0; i < 10; i++ {
		msg := togomq.NewMessage("events", []byte(fmt.Sprintf("event-%d", i)))
		msgChan <- msg
		log.Printf("Sent event-%d to channel\n", i)
		time.Sleep(100 * time.Millisecond) // Simulate some processing time
	}

	// Close channel to signal end of stream
	close(msgChan)

	// Wait for completion
	success := <-done
	if success {
		log.Println("Streaming publish completed successfully")
	}
}
