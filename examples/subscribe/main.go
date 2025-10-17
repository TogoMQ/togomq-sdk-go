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
		togomq.WithToken("your-token-here"),
		togomq.WithLogLevel("info"),
	)

	// Create client
	client, err := togomq.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Choose which example to run:
	// Uncomment the one you want to test

	// Example 1: Basic subscription to a specific topic
	// basicSubscribe(client)

	// Example 2: Subscribe to all topics with wildcard
	// subscribeAllTopics(client)

	// Example 3: Subscribe with pattern matching
	// subscribeWithPattern(client)

	// Example 4: Advanced subscription with options
	advancedSubscribe(client)
}

// basicSubscribe demonstrates basic subscription to a specific topic
func basicSubscribe(client *togomq.Client) {
	fmt.Println("=== Basic Subscription Example ===")

	// Subscribe to specific topic
	opts := togomq.NewSubscribeOptions("orders")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msgChan, errChan, err := client.Sub(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	log.Println("Listening for messages on topic 'orders'...")

	// Process messages
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				log.Println("Subscription ended")
				return
			}
			log.Printf("Received message from %s: %s\n", msg.Topic, string(msg.Body))
			log.Printf("Message UUID: %s\n", msg.UUID)

			// Access variables if present
			if len(msg.Variables) > 0 {
				log.Printf("Variables: %+v\n", msg.Variables)
				if priority, ok := msg.Variables["priority"]; ok {
					log.Printf("Priority: %s\n", priority)
				}
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

// subscribeAllTopics demonstrates subscribing to all topics using "*" wildcard
func subscribeAllTopics(client *togomq.Client) {
	fmt.Println("=== Subscribe to All Topics Example ===")

	// Subscribe to all topics using "*" wildcard
	opts := togomq.NewSubscribeOptions("*")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msgChan, errChan, err := client.Sub(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	log.Println("Listening for messages from ALL topics...")

	// Process messages
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				log.Println("Subscription ended")
				return
			}
			log.Printf("Received message from topic '%s': %s\n", msg.Topic, string(msg.Body))

		case err := <-errChan:
			log.Printf("Subscription error: %v\n", err)
			return

		case <-ctx.Done():
			log.Println("Context cancelled, ending subscription")
			return
		}
	}
}

// subscribeWithPattern demonstrates pattern matching with wildcards
func subscribeWithPattern(client *togomq.Client) {
	fmt.Println("=== Subscribe with Pattern Matching Example ===")

	// Subscribe to all orders topics (orders.new, orders.updated, etc.)
	opts := togomq.NewSubscribeOptions("orders.*")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msgChan, errChan, err := client.Sub(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	log.Println("Listening for messages matching pattern 'orders.*'...")

	// Process messages
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				log.Println("Subscription ended")
				return
			}
			log.Printf("Received message from %s: %s\n", msg.Topic, string(msg.Body))

		case err := <-errChan:
			log.Printf("Subscription error: %v\n", err)
			return

		case <-ctx.Done():
			log.Println("Context cancelled, ending subscription")
			return
		}
	}
}

// advancedSubscribe demonstrates subscription with batch size and rate limiting
func advancedSubscribe(client *togomq.Client) {
	fmt.Println("=== Advanced Subscription with Options ===")

	// Subscribe with batch size and rate limiting
	opts := togomq.NewSubscribeOptions("events").
		WithBatch(10).       // Receive up to 10 messages at once
		WithSpeedPerSec(100) // Limit to 100 messages per second

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msgChan, errChan, err := client.Sub(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	log.Println("Listening for messages with batch=10, speed=100/sec...")

	messageCount := 0

	// Process messages
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				log.Printf("Subscription ended. Total messages received: %d\n", messageCount)
				return
			}
			messageCount++
			log.Printf("[%d] Received: %s from %s\n", messageCount, string(msg.Body), msg.Topic)

		case err := <-errChan:
			log.Printf("Subscription error: %v\n", err)
			log.Printf("Messages received before error: %d\n", messageCount)
			return

		case <-ctx.Done():
			log.Printf("Context cancelled. Total messages received: %d\n", messageCount)
			return
		}
	}
}
