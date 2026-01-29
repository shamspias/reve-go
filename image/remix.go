package image

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/shamspias/reve-go/internal/transport"
	"github.com/shamspias/reve-go/internal/validator"
	"github.com/shamspias/reve-go/types"
)

// RemixParams contains parameters for remixing images.
type RemixParams struct {
	// Prompt describes the desired output (required).
	// Use <img>N</img> tags to reference images.
	// Maximum length: 2560 characters.
	Prompt string `json:"prompt"`

	// ReferenceImages is a list of base64 encoded images (required).
	// Minimum: 1, Maximum: 6
	ReferenceImages []string `json:"reference_images"`

	// AspectRatio is the desired aspect ratio.
	// Default: chosen by model
	AspectRatio types.AspectRatio `json:"aspect_ratio,omitempty"`

	// Version is the model version.
	// Default: latest
	Version types.ModelVersion `json:"version,omitempty"`

	// Postprocess contains postprocessing operations.
	Postprocess []types.Postprocess `json:"postprocessing,omitempty"`

	// TestTimeScaling controls quality (1-15).
	// Default: 1
	TestTimeScaling float64 `json:"test_time_scaling,omitempty"`

	// Breadcrumb is an optional tracking ID.
	Breadcrumb string `json:"-"`
}

// Validate validates the parameters.
func (p *RemixParams) Validate() error {
	if err := validator.ValidatePrompt(p.Prompt); err != nil {
		return err
	}
	if err := validator.ValidateReferenceImages(p.ReferenceImages); err != nil {
		return err
	}
	if err := validator.ValidateAspectRatio(string(p.AspectRatio)); err != nil {
		return err
	}
	if err := validator.ValidateScaling(p.TestTimeScaling); err != nil {
		return err
	}
	for _, pp := range p.Postprocess {
		if err := pp.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Remix combines multiple images with a text prompt.
//
// Example:
//
//	style, _ := types.NewImageFromFile("style.png")
//	content, _ := types.NewImageFromFile("content.png")
//
//	result, err := client.Images.Remix(ctx, &image.RemixParams{
//		Prompt: fmt.Sprintf("Apply style from %s to %s", types.Ref(0), types.Ref(1)),
//		ReferenceImages: []string{style.Base64(), content.Base64()},
//	})
func (s *Service) Remix(ctx context.Context, params *RemixParams) (*types.Result, error) {
	if params == nil {
		return nil, validator.ErrEmptyPrompt
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(ctx, &transport.Request{
		Method:     http.MethodPost,
		Path:       "/v1/image/remix",
		Body:       params,
		Breadcrumb: params.Breadcrumb,
	})
	if err != nil {
		return nil, err
	}

	var result types.Result
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// RemixRaw combines images and returns raw bytes.
//
// Example:
//
//	images := []string{img1.Base64(), img2.Base64()}
//	result, err := client.Images.RemixRaw(ctx, &image.RemixParams{
//		Prompt:          "Blend these styles",
//		ReferenceImages: images,
//		Version:         types.VersionLatestFast,
//	}, types.FormatWebP)
func (s *Service) RemixRaw(ctx context.Context, params *RemixParams, format types.OutputFormat) (*types.RawResult, error) {
	if params == nil {
		return nil, validator.ErrEmptyPrompt
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}

	if format == "" || format == types.FormatJSON {
		format = types.FormatPNG
	}

	resp, err := s.transport.DoRaw(ctx, &transport.Request{
		Method:     http.MethodPost,
		Path:       "/v1/image/remix",
		Body:       params,
		Accept:     string(format),
		Breadcrumb: params.Breadcrumb,
	})
	if err != nil {
		return nil, err
	}

	return &types.RawResult{
		Data:             resp.Data,
		ContentType:      resp.ContentType,
		Version:          resp.Version,
		ContentViolation: resp.ContentViolation,
		RequestID:        resp.RequestID,
		CreditsUsed:      resp.CreditsUsed,
		CreditsRemaining: resp.CreditsRemaining,
	}, nil
}
