package main

import (
	"context"
	"fmt"
	"log"

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

	ctx := context.Background()

	// Example 1: Count messages in a specific topic
	fmt.Println("=== Example 1: Count messages in specific topic ===")
	count, err := client.CountMessages(ctx, "orders")
	if err != nil {
		log.Fatalf("Failed to count messages: %v", err)
	}
	fmt.Printf("Messages in 'orders' topic: %d\n", count)

	// Example 2: Count messages using wildcard pattern
	fmt.Println("\n=== Example 2: Count messages with wildcard ===")
	count, err = client.CountMessages(ctx, "orders.*")
	if err != nil {
		log.Fatalf("Failed to count messages: %v", err)
	}
	fmt.Printf("Messages in 'orders.*' topics: %d\n", count)

	// Example 3: Count all messages across all topics
	fmt.Println("\n=== Example 3: Count all messages ===")
	count, err = client.CountMessages(ctx, "*")
	if err != nil {
		log.Fatalf("Failed to count messages: %v", err)
	}
	fmt.Printf("Total messages across all topics: %d\n", count)

	// Example 4: Count messages in multiple topics
	fmt.Println("\n=== Example 4: Count messages in multiple topics ===")
	topics := []string{"orders", "events", "notifications", "logs.*"}
	for _, topic := range topics {
		count, err := client.CountMessages(ctx, topic)
		if err != nil {
			log.Printf("Failed to count messages for topic '%s': %v", topic, err)
			continue
		}
		fmt.Printf("Topic %-15s: %d messages\n", topic, count)
	}
}
