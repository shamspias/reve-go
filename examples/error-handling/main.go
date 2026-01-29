// Example: Error Handling
//
// This example demonstrates comprehensive error handling.
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
	"strings"

	reve "github.com/shamspias/reve-go"
	"github.com/shamspias/reve-go/internal/transport"
	"github.com/shamspias/reve-go/internal/validator"
)

func main() {
	apiKey := os.Getenv("REVE_API_KEY")
	if apiKey == "" {
		log.Fatal("REVE_API_KEY required")
	}

	client := reve.NewClient(apiKey)
	ctx := context.Background()

	// =========================================================================
	// Example 1: Validation Errors
	// =========================================================================
	fmt.Println("=== Example 1: Validation Errors ===")

	// Empty prompt
	_, err := client.Images.Create(ctx, &reve.CreateParams{
		Prompt: "",
	})
	handleError("Empty prompt", err)

	// Prompt too long
	_, err = client.Images.Create(ctx, &reve.CreateParams{
		Prompt: strings.Repeat("a", 3000), // > 2560 chars
	})
	handleError("Prompt too long", err)

	// Invalid aspect ratio
	_, err = client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "test",
		AspectRatio: "invalid",
	})
	handleError("Invalid aspect ratio", err)

	// Invalid upscale factor
	_, err = client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "test",
		Postprocess: []reve.Postprocess{{Process: "upscale", UpscaleFactor: 10}},
	})
	handleError("Invalid upscale", err)

	// Invalid scaling
	_, err = client.Images.Create(ctx, &reve.CreateParams{
		Prompt:          "test",
		TestTimeScaling: 20, // > 15
	})
	handleError("Invalid scaling", err)

	// =========================================================================
	// Example 2: API Errors
	// =========================================================================
	fmt.Println("\n=== Example 2: API Errors ===")

	// Invalid API key
	badClient := reve.NewClient("invalid-key", reve.WithNoRetry())
	_, err = badClient.Images.Create(ctx, &reve.CreateParams{
		Prompt: "test",
	})
	handleError("Invalid API key", err)

	// =========================================================================
	// Example 3: Error Type Checking
	// =========================================================================
	fmt.Println("\n=== Example 3: Error Type Checking ===")

	// Check for specific validation errors
	_, err = client.Images.Create(ctx, &reve.CreateParams{Prompt: ""})

	if errors.Is(err, validator.ErrEmptyPrompt) {
		fmt.Println("✓ Detected: Empty prompt error")
	}

	// Check for API errors
	_, err = badClient.Images.Create(ctx, &reve.CreateParams{Prompt: "test"})

	var apiErr *transport.APIError
	if errors.As(err, &apiErr) {
		fmt.Println("✓ Detected: API error")
		fmt.Printf("  Code: %s\n", apiErr.Code)
		fmt.Printf("  Status: %d\n", apiErr.StatusCode)
		fmt.Printf("  Message: %s\n", apiErr.Message)

		// Check specific error conditions
		if apiErr.IsAuthError() {
			fmt.Println("  → Authentication failed")
		}
		if apiErr.IsRateLimit() {
			fmt.Println("  → Rate limited")
		}
		if apiErr.IsInsufficientFunds() {
			fmt.Println("  → Need more credits")
		}
		if apiErr.IsContentViolation() {
			fmt.Println("  → Content policy violation")
		}
		if apiErr.Retryable() {
			fmt.Println("  → Can be retried")
		}
	}

	// =========================================================================
	// Example 4: Recommended Error Handling Pattern
	// =========================================================================
	fmt.Println("\n=== Example 4: Recommended Pattern ===")

	result, err := client.Images.Create(ctx, &reve.CreateParams{
		Prompt: "A beautiful sunset",
	})

	if err != nil {
		// First, check validation errors
		switch {
		case errors.Is(err, validator.ErrEmptyPrompt):
			fmt.Println("Error: Please provide a prompt")
			return

		case errors.Is(err, validator.ErrPromptTooLong):
			fmt.Println("Error: Prompt is too long (max 2560 chars)")
			return

		case errors.Is(err, validator.ErrInvalidAspectRatio):
			fmt.Println("Error: Invalid aspect ratio")
			return

		case errors.Is(err, validator.ErrInvalidUpscaleFactor):
			fmt.Println("Error: Upscale factor must be 2, 3, or 4")
			return

		case errors.Is(err, validator.ErrInvalidScaling):
			fmt.Println("Error: Test time scaling must be 1-15")
			return
		}

		// Check API errors
		var apiErr *transport.APIError
		if errors.As(err, &apiErr) {
			switch {
			case apiErr.IsAuthError():
				fmt.Println("Error: Invalid API key")
				return

			case apiErr.IsRateLimit():
				fmt.Println("Error: Rate limited. Please wait and retry.")
				return

			case apiErr.IsInsufficientFunds():
				fmt.Println("Error: Insufficient credits. Please top up.")
				return

			case apiErr.IsContentViolation():
				fmt.Println("Error: Content policy violation. Modify your prompt.")
				return

			default:
				fmt.Printf("API Error: %s\n", apiErr.Message)
				if apiErr.Retryable() {
					fmt.Println("(This error can be retried)")
				}
				return
			}
		}

		// Check request errors (network, etc.)
		var reqErr *transport.RequestError
		if errors.As(err, &reqErr) {
			fmt.Printf("Request Error: %s failed: %v\n", reqErr.Op, reqErr.Err)
			return
		}

		// Unknown error
		fmt.Printf("Unknown error: %v\n", err)
		return
	}

	// Success!
	fmt.Printf("Success! Credits used: %d\n", result.CreditsUsed)
	result.SaveTo("sunset.png")

	// =========================================================================
	// Example 5: Edit-specific validation
	// =========================================================================
	fmt.Println("\n=== Example 5: Edit Validation ===")

	// Missing instruction
	_, err = client.Images.Edit(ctx, &reve.EditParams{
		Instruction:    "",
		ReferenceImage: "some-base64",
	})
	handleError("Empty instruction", err)

	// Missing reference image
	_, err = client.Images.Edit(ctx, &reve.EditParams{
		Instruction:    "Make it blue",
		ReferenceImage: "",
	})
	handleError("Empty reference image", err)

	// =========================================================================
	// Example 6: Remix-specific validation
	// =========================================================================
	fmt.Println("\n=== Example 6: Remix Validation ===")

	// No reference images
	_, err = client.Images.Remix(ctx, &reve.RemixParams{
		Prompt:          "Combine these",
		ReferenceImages: []string{},
	})
	handleError("No reference images", err)

	// Too many reference images
	_, err = client.Images.Remix(ctx, &reve.RemixParams{
		Prompt:          "Combine these",
		ReferenceImages: make([]string, 7), // > 6
	})
	handleError("Too many images", err)
}

func handleError(scenario string, err error) {
	if err != nil {
		fmt.Printf("%-25s → Error: %v\n", scenario, truncateError(err))
	} else {
		fmt.Printf("%-25s → Success (unexpected)\n", scenario)
	}
}

func truncateError(err error) string {
	msg := err.Error()
	if len(msg) > 60 {
		return msg[:60] + "..."
	}
	return msg
}
