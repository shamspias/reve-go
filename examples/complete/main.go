// Example: Complete Demo
//
// This example demonstrates all SDK features in one place.
//
// Run with:
//
//	REVE_API_KEY=your-key go run main.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	reve "github.com/shamspias/reve-go"
	"github.com/shamspias/reve-go/image"
	"github.com/shamspias/reve-go/internal/transport"
)

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║           Reve Go SDK - Complete Feature Demo                  ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	apiKey := os.Getenv("REVE_API_KEY")
	if apiKey == "" {
		log.Fatal("REVE_API_KEY environment variable is required")
	}

	// Run all demos
	demoClientConfiguration(apiKey)
	demoCreateImages(apiKey)
	demoEditImages(apiKey)
	demoRemixImages(apiKey)
	demoBatchOperations(apiKey)
	demoErrorHandling(apiKey)
	demoCostEstimation()

	fmt.Println("\n✅ All demos completed!")
}

func demoClientConfiguration(apiKey string) {
	fmt.Println("\n" + separator("Client Configuration"))

	// Basic client
	fmt.Println("1. Basic client:")
	client := reve.NewClient(apiKey)
	fmt.Printf("   Created with timeout=%v, retries=%d\n",
		client.Config().Timeout, client.Config().MaxRetries)

	// With options
	fmt.Println("\n2. With options:")
	client = reve.NewClient(apiKey,
		reve.WithTimeout(60*time.Second),
		reve.WithRetry(5, time.Second, 30*time.Second),
		reve.WithUserAgent("MyApp/1.0"),
		reve.WithDebug(false),
	)
	fmt.Printf("   Timeout=%v, Retries=%d, UserAgent=%s\n",
		client.Config().Timeout, client.Config().MaxRetries, client.Config().UserAgent)

	// Proxy configurations (examples only)
	fmt.Println("\n3. Proxy configurations:")
	fmt.Println("   HTTP:   reve.WithHTTPProxy(\"http://proxy:8080\")")
	fmt.Println("   SOCKS5: reve.WithSOCKS5Proxy(\"127.0.0.1:1080\", \"\", \"\")")
	fmt.Println("   Env:    reve.WithProxyFromEnvironment()")
}

func demoCreateImages(apiKey string) {
	fmt.Println("\n" + separator("Create Images"))

	client := reve.NewClient(apiKey)
	ctx := context.Background()

	// Basic create
	fmt.Println("1. Basic create:")
	result, err := client.Images.Create(ctx, &reve.CreateParams{
		Prompt: "A serene mountain lake at sunrise",
	})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Created! Credits: %d, Version: %s\n", result.CreditsUsed, result.Version)
		result.SaveTo("demo_basic.png")
	}

	// With options
	fmt.Println("\n2. With options (16:9, upscale):")
	result, err = client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "A cyberpunk city at night",
		AspectRatio: reve.Ratio16x9,
		Postprocess: []reve.Postprocess{reve.Upscale(2)},
	})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Created! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("demo_options.png")
	}
}

func demoEditImages(apiKey string) {
	fmt.Println("\n" + separator("Edit Images"))

	client := reve.NewClient(apiKey)
	ctx := context.Background()

	// First create an image to edit
	fmt.Println("1. Creating base image...")
	baseResult, err := client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "A red apple on white background",
		AspectRatio: reve.Ratio1x1,
	})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
		return
	}
	baseResult.SaveTo("demo_base.png")
	fmt.Printf("   ✓ Base image created\n")

	// Edit the image
	fmt.Println("\n2. Editing image (watercolor style):")
	result, err := client.Images.Edit(ctx, &reve.EditParams{
		Instruction:    "Convert to watercolor painting",
		ReferenceImage: baseResult.Image, // Use base64 directly
	})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Edited! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("demo_edit.png")
	}

	// Fast edit
	fmt.Println("\n3. Fast edit (5 credits):")
	result, err = client.Images.Edit(ctx, &reve.EditParams{
		Instruction:    "Add vintage filter",
		ReferenceImage: baseResult.Image,
		Version:        reve.VersionLatestFast,
	})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Fast edited! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("demo_fast_edit.png")
	}
}

func demoRemixImages(apiKey string) {
	fmt.Println("\n" + separator("Remix Images"))

	client := reve.NewClient(apiKey)
	ctx := context.Background()

	// Create two images to remix
	fmt.Println("1. Creating reference images...")

	style, err := client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "Abstract colorful geometric art",
		AspectRatio: reve.Ratio1x1,
	})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
		return
	}
	style.SaveTo("demo_style.png")

	content, err := client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "A peaceful countryside landscape",
		AspectRatio: reve.Ratio1x1,
	})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
		return
	}
	content.SaveTo("demo_content.png")
	fmt.Printf("   ✓ Reference images created\n")

	// Remix
	fmt.Println("\n2. Remixing images:")
	result, err := client.Images.Remix(ctx, &reve.RemixParams{
		Prompt: fmt.Sprintf("Apply style from %s to %s", reve.Ref(0), reve.Ref(1)),
		ReferenceImages: []string{
			style.Image,
			content.Image,
		},
	})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Remixed! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("demo_remix.png")
	}
}

func demoBatchOperations(apiKey string) {
	fmt.Println("\n" + separator("Batch Operations"))

	client := reve.NewClient(apiKey)
	ctx := context.Background()

	requests := []*reve.CreateParams{
		{Prompt: "A red rose"},
		{Prompt: "A white lily"},
		{Prompt: "A yellow sunflower"},
	}

	fmt.Printf("Processing %d requests concurrently...\n", len(requests))

	start := time.Now()
	results := client.Images.BatchCreate(ctx, requests, &image.BatchConfig{
		Concurrency: 3,
	})
	elapsed := time.Since(start)

	fmt.Printf("   Completed in %v\n", elapsed)
	fmt.Printf("   Success: %d, Failed: %d\n",
		image.SuccessCount(results), image.ErrorCount(results))

	totalCredits := 0
	for _, r := range results {
		if r.Result != nil {
			totalCredits += r.Result.CreditsUsed
			r.Result.SaveTo(fmt.Sprintf("demo_batch_%d.png", r.Index))
		}
	}
	fmt.Printf("   Total credits: %d\n", totalCredits)
}

func demoErrorHandling(apiKey string) {
	fmt.Println("\n" + separator("Error Handling"))

	client := reve.NewClient(apiKey)
	ctx := context.Background()

	// Validation error
	fmt.Println("1. Validation error (empty prompt):")
	_, err := client.Images.Create(ctx, &reve.CreateParams{Prompt: ""})
	if err != nil {
		fmt.Printf("   ✓ Caught: %v\n", err)
	}

	// API error (invalid key)
	fmt.Println("\n2. API error (invalid key):")
	badClient := reve.NewClient("invalid-key", reve.WithNoRetry())
	_, err = badClient.Images.Create(ctx, &reve.CreateParams{Prompt: "test"})
	if err != nil {
		var apiErr *transport.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("   ✓ Caught API error: %s (status=%d)\n", apiErr.Code, apiErr.StatusCode)
			fmt.Printf("   IsAuthError: %v, Retryable: %v\n", apiErr.IsAuthError(), apiErr.Retryable())
		}
	}
}

func demoCostEstimation() {
	fmt.Println("\n" + separator("Cost Estimation"))

	fmt.Println("Create costs:")
	fmt.Printf("  Basic:          %s\n", reve.EstimateCreate(1, nil))
	fmt.Printf("  With 2x scale:  %s\n", reve.EstimateCreate(2, nil))
	fmt.Printf("  With upscale:   %s\n", reve.EstimateCreate(1, []reve.Postprocess{reve.Upscale(2)}))

	fmt.Println("\nEdit costs:")
	fmt.Printf("  Standard:       %s\n", reve.EstimateEdit(false, 1, nil))
	fmt.Printf("  Fast:           %s\n", reve.EstimateEdit(true, 1, nil))

	fmt.Println("\nRemix costs:")
	fmt.Printf("  Standard:       %s\n", reve.EstimateRemix(false, 1, nil))
	fmt.Printf("  Fast:           %s\n", reve.EstimateRemix(true, 1, nil))
}

func separator(title string) string {
	return fmt.Sprintf("═══ %s ═══", title)
}
