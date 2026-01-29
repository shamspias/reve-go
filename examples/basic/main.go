// Example: Basic Usage
//
// This example demonstrates the simplest way to use the Reve SDK.
//
// Run with:
//
//	REVE_API_KEY=your-key go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	reve "github.com/shamspias/reve-go"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("REVE_API_KEY")
	if apiKey == "" {
		log.Fatal("REVE_API_KEY environment variable is required")
	}

	// Create client
	client := reve.NewClient(apiKey)

	// Generate an image
	result, err := client.Images.Create(context.Background(), &reve.CreateParams{
		Prompt: "A beautiful mountain landscape at sunset with golden light",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print result info
	fmt.Println("Image generated successfully!")
	fmt.Printf("  Request ID: %s\n", result.RequestID)
	fmt.Printf("  Version: %s\n", result.Version)
	fmt.Printf("  Credits Used: %d\n", result.CreditsUsed)
	fmt.Printf("  Credits Remaining: %d\n", result.CreditsRemaining)

	// Save the image
	if err := result.SaveTo("output.png"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Saved to: output.png")
}
