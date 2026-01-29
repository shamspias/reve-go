// Package image provides the image generation service for the Reve API.
//
// This package handles all image-related operations:
//   - Create: Generate images from text descriptions
//   - Edit: Modify images with text instructions
//   - Remix: Combine images with text prompts
//   - Batch operations for concurrent processing
//
// # Usage
//
//	import (
//		reve "github.com/shamspias/reve-go"
//		"github.com/shamspias/reve-go/image"
//	)
//
//	client := reve.NewClient(apiKey)
//
//	result, err := client.Images.Create(ctx, &image.CreateParams{
//		Prompt: "A beautiful sunset",
//	})
package image

import (
	"github.com/shamspias/reve-go/internal/transport"
)

// Service handles image operations.
type Service struct {
	transport *transport.Client
}

// NewService creates a new image service.
func NewService(t *transport.Client) *Service {
	return &Service{transport: t}
}
