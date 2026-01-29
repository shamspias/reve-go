// Example: Create Images
//
// This example demonstrates all image creation options.
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
)

func main() {
	apiKey := os.Getenv("REVE_API_KEY")
	if apiKey == "" {
		log.Fatal("REVE_API_KEY required")
	}

	client := reve.NewClient(apiKey,
		reve.WithTimeout(60*time.Second),
		reve.WithDebug(true),
	)

	ctx := context.Background()

	// =========================================================================
	// Example 1: Basic create
	// =========================================================================
	fmt.Println("=== Example 1: Basic Create ===")
	result, err := client.Images.Create(ctx, &reve.CreateParams{
		Prompt: "A serene Japanese garden with cherry blossoms",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Created! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("basic.png")
	}

	// =========================================================================
	// Example 2: With aspect ratio
	// =========================================================================
	fmt.Println("\n=== Example 2: With Aspect Ratio ===")
	result, err = client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "A cyberpunk cityscape at night with neon lights",
		AspectRatio: reve.Ratio16x9, // Widescreen
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Created 16:9 image! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("widescreen.png")
	}

	// =========================================================================
	// Example 3: Square image
	// =========================================================================
	fmt.Println("\n=== Example 3: Square Image ===")
	result, err = client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "A minimalist logo design",
		AspectRatio: reve.Ratio1x1,
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Created square image! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("square.png")
	}

	// =========================================================================
	// Example 4: Portrait for mobile
	// =========================================================================
	fmt.Println("\n=== Example 4: Portrait (Mobile) ===")
	result, err = client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "A tall waterfall in a misty forest",
		AspectRatio: reve.Ratio9x16, // Mobile portrait
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Created portrait image! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("portrait.png")
	}

	// =========================================================================
	// Example 5: With test time scaling (better quality)
	// =========================================================================
	fmt.Println("\n=== Example 5: High Quality (Test Time Scaling) ===")
	result, err = client.Images.Create(ctx, &reve.CreateParams{
		Prompt:          "A photorealistic portrait of a majestic lion",
		AspectRatio:     reve.Ratio1x1,
		TestTimeScaling: 3, // 3x quality
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Created high-quality image! Credits: %d (scaled)\n", result.CreditsUsed)
		result.SaveTo("highquality.png")
	}

	// =========================================================================
	// Example 6: With upscale postprocessing
	// =========================================================================
	fmt.Println("\n=== Example 6: With Upscale ===")
	result, err = client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "A detailed fantasy castle on a cliff",
		AspectRatio: reve.Ratio16x9,
		Postprocess: []reve.Postprocess{
			reve.Upscale(2), // 2x upscale
		},
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Created upscaled image! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("upscaled.png")
	}

	// =========================================================================
	// Example 7: With background removal
	// =========================================================================
	fmt.Println("\n=== Example 7: Background Removal ===")
	result, err = client.Images.Create(ctx, &reve.CreateParams{
		Prompt:      "A red sports car, product photography",
		AspectRatio: reve.Ratio4x3,
		Postprocess: []reve.Postprocess{
			reve.RemoveBackground(),
		},
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Created with transparent background! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("transparent.png")
	}

	// =========================================================================
	// Example 8: Full options
	// =========================================================================
	fmt.Println("\n=== Example 8: All Options ===")
	result, err = client.Images.Create(ctx, &reve.CreateParams{
		Prompt:          "An astronaut riding a horse on Mars, cinematic",
		AspectRatio:     reve.Ratio16x9,
		Version:         reve.VersionLatest,
		TestTimeScaling: 2,
		Postprocess: []reve.Postprocess{
			reve.Upscale(2),
		},
		Breadcrumb: "astronaut-mars-v1", // Tracking ID
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Created with all options! Credits: %d\n", result.CreditsUsed)
		result.SaveTo("full_options.png")
	}

	// =========================================================================
	// Example 9: Get raw bytes (different format)
	// =========================================================================
	fmt.Println("\n=== Example 9: Raw JPEG Response ===")
	rawResult, err := client.Images.CreateRaw(ctx, &reve.CreateParams{
		Prompt:      "A sunset over the ocean",
		AspectRatio: reve.Ratio3x2,
	}, reve.FormatJPEG)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Received %d bytes as %s\n", rawResult.Size(), rawResult.ContentType)
		rawResult.SaveTo("sunset.jpg")
	}

	// =========================================================================
	// Cost estimation
	// =========================================================================
	fmt.Println("\n=== Cost Estimation ===")
	fmt.Printf("Basic create: %s\n", reve.EstimateCreate(1, nil))
	fmt.Printf("With 2x scaling: %s\n", reve.EstimateCreate(2, nil))
	fmt.Printf("With upscale: %s\n", reve.EstimateCreate(1, []reve.Postprocess{reve.Upscale(2)}))
	fmt.Printf("Full options: %s\n", reve.EstimateCreate(2, []reve.Postprocess{reve.Upscale(2)}))
}
