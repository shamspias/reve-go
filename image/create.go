package image

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/shamspias/reve-go/internal/transport"
	"github.com/shamspias/reve-go/internal/validator"
	"github.com/shamspias/reve-go/types"
)

// CreateParams contains parameters for creating an image.
type CreateParams struct {
	// Prompt is the text description (required).
	// Maximum length: 2560 characters.
	Prompt string `json:"prompt"`

	// AspectRatio is the desired aspect ratio.
	// Default: 3:2
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
func (p *CreateParams) Validate() error {
	if err := validator.ValidatePrompt(p.Prompt); err != nil {
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

// Create generates an image from a text description.
//
// Example:
//
//	result, err := client.Images.Create(ctx, &image.CreateParams{
//		Prompt:      "A serene mountain lake at sunrise",
//		AspectRatio: types.Ratio16x9,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = result.SaveTo("lake.png")
func (s *Service) Create(ctx context.Context, params *CreateParams) (*types.Result, error) {
	if params == nil {
		return nil, validator.ErrEmptyPrompt
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(ctx, &transport.Request{
		Method:     http.MethodPost,
		Path:       "/v1/image/create",
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

// CreateRaw generates an image and returns raw bytes.
//
// Example:
//
//	result, err := client.Images.CreateRaw(ctx, &image.CreateParams{
//		Prompt: "A sunset over the ocean",
//	}, types.FormatPNG)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = result.SaveTo("sunset.png")
func (s *Service) CreateRaw(ctx context.Context, params *CreateParams, format types.OutputFormat) (*types.RawResult, error) {
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
		Path:       "/v1/image/create",
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
