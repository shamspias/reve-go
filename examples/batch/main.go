// Example: Batch Operations
//
// This example demonstrates concurrent batch processing.
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
	"time"

	reve "github.com/shamspias/reve-go"
	"github.com/shamspias/reve-go/image"
)

func main() {
	apiKey := os.Getenv("REVE_API_KEY")
	if apiKey == "" {
		log.Fatal("REVE_API_KEY required")
	}

	client := reve.NewClient(apiKey)
	ctx := context.Background()

	// =========================================================================
	// Example 1: Basic batch create
	// =========================================================================
	fmt.Println("=== Example 1: Basic Batch Create ===")

	requests := []*reve.CreateParams{
		{Prompt: "A red apple on a wooden table"},
		{Prompt: "A green pear in a ceramic bowl"},
		{Prompt: "A bunch of purple grapes on a vine"},
		{Prompt: "A sliced orange showing segments"},
		{Prompt: "A yellow banana on a white plate"},
	}

	start := time.Now()
	results := client.Images.BatchCreate(ctx, requests, nil)
	elapsed := time.Since(start)

	fmt.Printf("Processed %d requests in %v\n", len(results), elapsed)
	fmt.Printf("Success: %d, Failed: %d\n",
		image.SuccessCount(results), image.ErrorCount(results))

	for _, r := range results {
		if r.Error != nil {
			fmt.Printf("  [%d] Failed: %v\n", r.Index, r.Error)
		} else {
			filename := fmt.Sprintf("fruit_%d.png", r.Index)
			r.Result.SaveTo(filename)
			fmt.Printf("  [%d] Success: %s (%d credits)\n", r.Index, filename, r.Result.CreditsUsed)
		}
	}

	// =========================================================================
	// Example 2: With concurrency control
	// =========================================================================
	fmt.Println("\n=== Example 2: Controlled Concurrency ===")

	requests2 := []*reve.CreateParams{
		{Prompt: "A modern skyscraper"},
		{Prompt: "A cozy cottage"},
		{Prompt: "A beach house"},
		{Prompt: "A mountain cabin"},
	}

	start = time.Now()
	results = client.Images.BatchCreate(ctx, requests2, &image.BatchConfig{
		Concurrency: 2, // Only 2 concurrent requests
		StopOnError: false,
	})
	elapsed = time.Since(start)

	fmt.Printf("Processed with concurrency=2 in %v\n", elapsed)
	fmt.Printf("Success: %d\n", image.SuccessCount(results))

	// =========================================================================
	// Example 3: Stop on first error
	// =========================================================================
	fmt.Println("\n=== Example 3: Stop on Error ===")

	requests3 := []*reve.CreateParams{
		{Prompt: "A valid image request"},
		{Prompt: ""}, // Invalid: empty prompt will fail
		{Prompt: "This won't run if StopOnError is true"},
	}

	results = client.Images.BatchCreate(ctx, requests3, &image.BatchConfig{
		Concurrency: 1,
		StopOnError: true, // Stop after first error
	})

	fmt.Printf("Processed: %d, Errors: %d\n",
		image.SuccessCount(results), image.ErrorCount(results))

	// =========================================================================
	// Example 4: Batch edit
	// =========================================================================
	fmt.Println("\n=== Example 4: Batch Edit ===")

	// First create a base image
	baseResult, err := client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "A simple red car",
		AspectRatio: reve.Ratio1x1,
	})
	if err != nil {
		log.Fatal(err)
	}

	baseImg, _ := reve.NewImageFromFile("") // Would load actual file
	baseB64 := baseResult.Image             // Use the base64 from result

	editRequests := []*reve.EditParams{
		{Instruction: "Change the car color to blue", ReferenceImage: baseB64},
		{Instruction: "Change the car color to green", ReferenceImage: baseB64},
		{Instruction: "Change the car color to yellow", ReferenceImage: baseB64},
	}

	_ = baseImg // Unused in this example

	editResults := client.Images.BatchEdit(ctx, editRequests, &image.BatchConfig{
		Concurrency: 3,
	})

	fmt.Printf("Batch edit: %d/%d successful\n",
		image.SuccessCount(editResults), len(editResults))

	// =========================================================================
	// Example 5: Using helper functions
	// =========================================================================
	fmt.Println("\n=== Example 5: Result Helpers ===")

	// Get only successful results
	successful := image.Successful(results)
	fmt.Printf("Got %d successful results\n", len(successful))

	// Get all errors
	errors := image.Errors(results)
	if len(errors) > 0 {
		fmt.Printf("Errors encountered:\n")
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
	}

	// =========================================================================
	// Example 6: Large batch with progress
	// =========================================================================
	fmt.Println("\n=== Example 6: Large Batch ===")

	largeBatch := make([]*reve.CreateParams, 10)
	for i := range largeBatch {
		largeBatch[i] = &reve.CreateParams{
			Prompt:      fmt.Sprintf("Abstract art piece number %d", i+1),
			AspectRatio: reve.Ratio1x1,
		}
	}

	fmt.Printf("Processing %d requests...\n", len(largeBatch))

	start = time.Now()
	largeResults := client.Images.BatchCreate(ctx, largeBatch, &image.BatchConfig{
		Concurrency: 5,
	})
	elapsed = time.Since(start)

	totalCredits := 0
	for _, r := range largeResults {
		if r.Result != nil {
			totalCredits += r.Result.CreditsUsed
		}
	}

	fmt.Printf("Completed in %v\n", elapsed)
	fmt.Printf("Success rate: %d/%d (%.1f%%)\n",
		image.SuccessCount(largeResults),
		len(largeResults),
		float64(image.SuccessCount(largeResults))/float64(len(largeResults))*100)
	fmt.Printf("Total credits used: %d\n", totalCredits)
}
