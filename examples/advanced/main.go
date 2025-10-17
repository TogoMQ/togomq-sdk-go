package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/TogoMQ/togomq-sdk-go"
)

func main() {
	// Create client configuration
	config := togomq.NewConfig(
		togomq.WithToken("your-token-here"),
		togomq.WithLogLevel("debug"),
	)

	// Create client
	client, err := togomq.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Choose which advanced example to run:
	// Uncomment the one you want to test

	// Example 1: Concurrent publishing and subscribing
	// concurrentPubSub(client)

	// Example 2: Error handling and retry logic
	// errorHandlingExample(client)

	// Example 3: Custom configuration and connection options
	customConfigExample()

	// Example 4: Message with full feature set
	// fullFeaturedMessages(client)
}

// concurrentPubSub demonstrates concurrent publishing and subscribing
func concurrentPubSub(client *togomq.Client) {
	fmt.Println("=== Concurrent Pub/Sub Example ===")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// Start subscriber in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		opts := togomq.NewSubscribeOptions("concurrent-test")
		msgChan, errChan, err := client.Sub(ctx, opts)
		if err != nil {
			log.Printf("Subscribe error: %v", err)
			return
		}

		log.Println("Subscriber started, waiting for messages...")
		count := 0

		for {
			select {
			case msg, ok := <-msgChan:
				if !ok {
					log.Printf("Subscriber finished. Received %d messages\n", count)
					return
				}
				count++
				log.Printf("Subscriber received message %d: %s\n", count, string(msg.Body))

			case err := <-errChan:
				log.Printf("Subscriber error: %v\n", err)
				return

			case <-ctx.Done():
				log.Printf("Subscriber context done. Received %d messages\n", count)
				return
			}
		}
	}()

	// Give subscriber time to connect
	time.Sleep(2 * time.Second)

	// Start multiple publishers in goroutines
	for i := 0; i < 3; i++ {
		wg.Add(1)
		publisherID := i
		go func() {
			defer wg.Done()

			for j := 0; j < 5; j++ {
				msg := togomq.NewMessage(
					"concurrent-test",
					[]byte(fmt.Sprintf("Message from publisher %d, msg %d", publisherID, j)),
				).WithVariables(map[string]string{
					"publisher": fmt.Sprintf("%d", publisherID),
					"sequence":  fmt.Sprintf("%d", j),
				})

				resp, err := client.PubBatch(ctx, []*togomq.Message{msg})
				if err != nil {
					log.Printf("Publisher %d error: %v\n", publisherID, err)
					return
				}
				log.Printf("Publisher %d sent message %d (received: %d)\n",
					publisherID, j, resp.MessagesReceived)

				time.Sleep(500 * time.Millisecond)
			}
			log.Printf("Publisher %d finished\n", publisherID)
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()
	log.Println("All concurrent operations completed")
}

// errorHandlingExample demonstrates comprehensive error handling
func errorHandlingExample(client *togomq.Client) {
	fmt.Println("=== Error Handling Example ===")

	ctx := context.Background()

	// Example 1: Handle validation errors
	invalidMessage := togomq.NewMessage("", []byte("no topic"))
	_, err := client.PubBatch(ctx, []*togomq.Message{invalidMessage})
	if err != nil {
		handleError(err, "Publishing invalid message")
	}

	// Example 2: Handle timeout errors
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
	defer cancel()

	msg := togomq.NewMessage("test-topic", []byte("test"))
	_, err = client.PubBatch(timeoutCtx, []*togomq.Message{msg})
	if err != nil {
		handleError(err, "Publishing with timeout")
	}

	// Example 3: Subscribe with invalid topic
	opts := togomq.NewSubscribeOptions("")
	_, _, err = client.Sub(ctx, opts)
	if err != nil {
		handleError(err, "Subscribing with invalid topic")
	}
}

// handleError demonstrates how to handle different types of TogoMQ errors
func handleError(err error, operation string) {
	log.Printf("Error during %s: %v\n", operation, err)

	// Check if it's a TogoMQ error
	if togomqErr, ok := err.(*togomq.TogoMQError); ok {
		switch togomqErr.Code {
		case togomq.ErrCodeAuth:
			log.Println("  → Authentication failed. Check your token.")
		case togomq.ErrCodeConnection:
			log.Println("  → Connection error. Check server availability.")
		case togomq.ErrCodeValidation:
			log.Println("  → Validation error. Check your input parameters.")
		case togomq.ErrCodePublish:
			log.Println("  → Publishing error. Check message format.")
		case togomq.ErrCodeSubscribe:
			log.Println("  → Subscription error. Check topic and options.")
		case togomq.ErrCodeConfiguration:
			log.Println("  → Configuration error. Check client setup.")
		default:
			log.Printf("  → Error code: %s\n", togomqErr.Code)
		}
	}
}

// customConfigExample demonstrates custom configuration options
func customConfigExample() {
	fmt.Println("=== Custom Configuration Example ===")

	// Example with custom host and port
	config := togomq.NewConfig(
		togomq.WithHost("custom.togomq.io"),
		togomq.WithPort(9000),
		togomq.WithLogLevel("debug"),
		togomq.WithToken("your-custom-token"),
	)

	client, err := togomq.NewClient(config)
	if err != nil {
		log.Printf("Failed to create client with custom config: %v", err)
		log.Println("  → This is expected when using a non-existent custom server")
		return
	}
	defer client.Close()

	log.Println("Successfully created client with custom configuration")

	// Try to use the client
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := togomq.NewMessage("test", []byte("test message"))
	_, err = client.PubBatch(ctx, []*togomq.Message{msg})
	if err != nil {
		handleError(err, "Using custom configured client")
	}
}

// fullFeaturedMessages demonstrates messages with all available features
func fullFeaturedMessages(client *togomq.Client) {
	fmt.Println("=== Full Featured Messages Example ===")

	ctx := context.Background()

	// Create messages with all features
	messages := []*togomq.Message{
		// Simple message
		togomq.NewMessage("notifications", []byte("Simple notification")),

		// Message with variables
		togomq.NewMessage("orders", []byte(`{"order_id": 12345, "total": 99.99}`)).
			WithVariables(map[string]string{
				"customer_id": "CUST-123",
				"priority":    "high",
				"region":      "us-west",
				"order_type":  "express",
			}),

		// Message with postpone (delayed delivery)
		togomq.NewMessage("scheduled-tasks", []byte("Execute this task later")).
			WithPostpone(300).
			WithVariables(map[string]string{
				"scheduled_by": "system",
				"task_type":    "cleanup",
			}),

		// Message with retention
		togomq.NewMessage("ephemeral-events", []byte("Short-lived event")).
			WithRetention(3600).
			WithVariables(map[string]string{
				"event_type": "user_action",
				"ttl":        "1h",
			}),

		// Message with both postpone and retention
		togomq.NewMessage("complex-workflow", []byte("Complex workflow step")).
			WithPostpone(60).
			WithRetention(7200).
			WithVariables(map[string]string{
				"workflow_id": "WF-456",
				"step":        "3",
				"retry_count": "0",
			}),
	}

	// Publish all messages
	resp, err := client.PubBatch(ctx, messages)
	if err != nil {
		handleError(err, "Publishing full-featured messages")
		return
	}

	log.Printf("Successfully published %d full-featured messages\n", resp.MessagesReceived)

	// Show details of what was published
	for i, msg := range messages {
		log.Printf("Message %d:\n", i+1)
		log.Printf("  Topic: %s\n", msg.Topic)
		log.Printf("  Body: %s\n", string(msg.Body))
		if msg.Postpone > 0 {
			log.Printf("  Postpone: %d seconds\n", msg.Postpone)
		}
		if msg.Retention > 0 {
			log.Printf("  Retention: %d seconds\n", msg.Retention)
		}
		if len(msg.Variables) > 0 {
			log.Printf("  Variables: %+v\n", msg.Variables)
		}
	}
}
