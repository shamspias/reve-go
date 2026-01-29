// Example: Remix Images
//
// This example demonstrates combining multiple images with text prompts.
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

	// First, create some test images
	fmt.Println("Creating test images...")

	// Style image
	styleResult, err := client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "Abstract colorful geometric art, vibrant colors",
		AspectRatio: reve.Ratio1x1,
	})
	if err != nil {
		log.Fatal(err)
	}
	styleResult.SaveTo("style.png")
	fmt.Println("Style image saved to style.png")

	// Content image
	contentResult, err := client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "A peaceful countryside landscape with rolling hills",
		AspectRatio: reve.Ratio1x1,
	})
	if err != nil {
		log.Fatal(err)
	}
	contentResult.SaveTo("content.png")
	fmt.Println("Content image saved to content.png")

	// Load images
	styleImg, _ := reve.NewImageFromFile("style.png")
	contentImg, _ := reve.NewImageFromFile("content.png")

	// =========================================================================
	// Example 1: Style transfer
	// =========================================================================
	fmt.Println("\n=== Example 1: Style Transfer ===")
	result, err := client.Images.Remix(ctx, &reve.RemixParams{
		Prompt: fmt.Sprintf("Apply the artistic style from %s to the landscape in %s",
			reve.Ref(0), reve.Ref(1)),
		ReferenceImages: []string{
			styleImg.Base64(),
			contentImg.Base64(),
		},
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Style transferred! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("style_transfer.png")
	}

	// =========================================================================
	// Example 2: Blend images
	// =========================================================================
	fmt.Println("\n=== Example 2: Blend Images ===")
	result, err = client.Images.Remix(ctx, &reve.RemixParams{
		Prompt: fmt.Sprintf("Blend the colors from %s with the composition of %s harmoniously",
			reve.Ref(0), reve.Ref(1)),
		ReferenceImages: []string{
			styleImg.Base64(),
			contentImg.Base64(),
		},
		AspectRatio: reve.Ratio16x9,
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Images blended! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("blended.png")
	}

	// =========================================================================
	// Example 3: Fast mode remix
	// =========================================================================
	fmt.Println("\n=== Example 3: Fast Mode ===")
	result, err = client.Images.Remix(ctx, &reve.RemixParams{
		Prompt: fmt.Sprintf("Create a dreamlike fusion of %s and %s",
			reve.Ref(0), reve.Ref(1)),
		ReferenceImages: []string{
			styleImg.Base64(),
			contentImg.Base64(),
		},
		Version: reve.VersionLatestFast, // 5 credits instead of 30
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Fast remix! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("fast_remix.png")
	}

	// =========================================================================
	// Example 4: Using single image as reference
	// =========================================================================
	fmt.Println("\n=== Example 4: Single Reference ===")
	result, err = client.Images.Remix(ctx, &reve.RemixParams{
		Prompt: fmt.Sprintf("Create a winter version of %s with snow", reve.Ref(0)),
		ReferenceImages: []string{
			contentImg.Base64(),
		},
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Winter version! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("winter.png")
	}

	// =========================================================================
	// Example 5: Multiple image references (3+ images)
	// =========================================================================
	fmt.Println("\n=== Example 5: Multiple References ===")

	// Create a third image
	thirdResult, _ := client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "Golden sunset lighting",
		AspectRatio: reve.Ratio1x1,
	})
	thirdResult.SaveTo("lighting.png")
	lightingImg, _ := reve.NewImageFromFile("lighting.png")

	result, err = client.Images.Remix(ctx, &reve.RemixParams{
		Prompt: fmt.Sprintf("Combine the style of %s, the scene from %s, and the lighting from %s",
			reve.Ref(0), reve.Ref(1), reve.Ref(2)),
		ReferenceImages: []string{
			styleImg.Base64(),
			contentImg.Base64(),
			lightingImg.Base64(),
		},
		AspectRatio: reve.Ratio16x9,
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Multi-reference remix! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("multi_remix.png")
	}

	// =========================================================================
	// Example 6: With postprocessing
	// =========================================================================
	fmt.Println("\n=== Example 6: With Upscale ===")
	result, err = client.Images.Remix(ctx, &reve.RemixParams{
		Prompt: fmt.Sprintf("Merge %s and %s into a surreal artwork",
			reve.Ref(0), reve.Ref(1)),
		ReferenceImages: []string{
			styleImg.Base64(),
			contentImg.Base64(),
		},
		Postprocess: []reve.Postprocess{
			reve.Upscale(2),
		},
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Upscaled remix! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("upscaled_remix.png")
	}

	// =========================================================================
	// Cost comparison
	// =========================================================================
	fmt.Println("\n=== Cost Comparison ===")
	fmt.Printf("Standard remix: %s\n", reve.EstimateRemix(false, 1, nil))
	fmt.Printf("Fast remix: %s\n", reve.EstimateRemix(true, 1, nil))
	fmt.Printf("With 2x upscale: %s\n", reve.EstimateRemix(false, 1, []reve.Postprocess{reve.Upscale(2)}))
}
