// Example: Edit Images
//
// This example demonstrates image editing capabilities.
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
	apiKey := os.Getenv("REVE_API_KEY")
	if apiKey == "" {
		log.Fatal("REVE_API_KEY required")
	}

	client := reve.NewClient(apiKey)
	ctx := context.Background()

	// First, create a test image to edit
	fmt.Println("Creating test image...")
	testResult, err := client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "A simple red apple on a white background",
		AspectRatio: reve.Ratio1x1,
	})
	if err != nil {
		log.Fatal(err)
	}
	testResult.SaveTo("original.png")
	fmt.Println("Original saved to original.png")

	// Load the image for editing
	img, err := reve.NewImageFromFile("original.png")
	if err != nil {
		log.Fatal(err)
	}

	// =========================================================================
	// Example 1: Style transfer
	// =========================================================================
	fmt.Println("\n=== Example 1: Watercolor Style ===")
	result, err := client.Images.Edit(ctx, &reve.EditParams{
		Instruction:    "Convert to a watercolor painting style",
		ReferenceImage: img.Base64(),
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Edited! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("watercolor.png")
	}

	// =========================================================================
	// Example 2: Color adjustment
	// =========================================================================
	fmt.Println("\n=== Example 2: Color Change ===")
	result, err = client.Images.Edit(ctx, &reve.EditParams{
		Instruction:    "Change the apple color from red to green",
		ReferenceImage: img.Base64(),
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Color changed! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("green_apple.png")
	}

	// =========================================================================
	// Example 3: Add elements
	// =========================================================================
	fmt.Println("\n=== Example 3: Add Elements ===")
	result, err = client.Images.Edit(ctx, &reve.EditParams{
		Instruction:    "Add water droplets on the apple surface",
		ReferenceImage: img.Base64(),
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Elements added! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("droplets.png")
	}

	// =========================================================================
	// Example 4: Background change
	// =========================================================================
	fmt.Println("\n=== Example 4: Change Background ===")
	result, err = client.Images.Edit(ctx, &reve.EditParams{
		Instruction:    "Change the background to a wooden table in a rustic kitchen",
		ReferenceImage: img.Base64(),
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Background changed! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("rustic.png")
	}

	// =========================================================================
	// Example 5: Fast mode (cheaper, faster)
	// =========================================================================
	fmt.Println("\n=== Example 5: Fast Mode ===")
	result, err = client.Images.Edit(ctx, &reve.EditParams{
		Instruction:    "Add a vintage sepia filter",
		ReferenceImage: img.Base64(),
		Version:        reve.VersionLatestFast, // Fast mode: 5 credits
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Fast edit! Credits: %d (only 5 for fast mode)\n", result.CreditsUsed)
		result.SaveTo("vintage.png")
	}

	// =========================================================================
	// Example 6: With aspect ratio change
	// =========================================================================
	fmt.Println("\n=== Example 6: Change Aspect Ratio ===")
	result, err = client.Images.Edit(ctx, &reve.EditParams{
		Instruction:    "Extend the image with more background",
		ReferenceImage: img.Base64(),
		AspectRatio:    reve.Ratio16x9, // Widescreen
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Aspect changed! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("widescreen.png")
	}

	// =========================================================================
	// Example 7: High quality with upscale
	// =========================================================================
	fmt.Println("\n=== Example 7: High Quality Edit ===")
	result, err = client.Images.Edit(ctx, &reve.EditParams{
		Instruction:     "Transform into an oil painting masterpiece",
		ReferenceImage:  img.Base64(),
		TestTimeScaling: 2,
		Postprocess: []reve.Postprocess{
			reve.Upscale(2),
		},
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("High quality! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("oilpainting.png")
	}

	// =========================================================================
	// Cost comparison
	// =========================================================================
	fmt.Println("\n=== Cost Comparison ===")
	fmt.Printf("Standard edit: %s\n", reve.EstimateEdit(false, 1, nil))
	fmt.Printf("Fast edit: %s\n", reve.EstimateEdit(true, 1, nil))
	fmt.Printf("With 2x scaling: %s\n", reve.EstimateEdit(false, 2, nil))
	fmt.Printf("Fast + upscale: %s\n", reve.EstimateEdit(true, 1, []reve.Postprocess{reve.Upscale(2)}))
}
