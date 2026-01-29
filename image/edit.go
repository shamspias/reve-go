package image

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/shamspias/reve-go/internal/transport"
	"github.com/shamspias/reve-go/internal/validator"
	"github.com/shamspias/reve-go/types"
)

// EditParams contains parameters for editing an image.
type EditParams struct {
	// Instruction describes how to edit the image (required).
	// Maximum length: 2560 characters.
	Instruction string `json:"edit_instruction"`

	// ReferenceImage is the base64 encoded image (required).
	ReferenceImage string `json:"reference_image"`

	// AspectRatio is the desired aspect ratio.
	// Default: reference image aspect ratio
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
func (p *EditParams) Validate() error {
	if err := validator.ValidateInstruction(p.Instruction); err != nil {
		return err
	}
	if err := validator.ValidateReferenceImage(p.ReferenceImage); err != nil {
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

// Edit modifies an image based on text instructions.
//
// Example:
//
//	img, _ := types.NewImageFromFile("photo.jpg")
//	result, err := client.Images.Edit(ctx, &image.EditParams{
//		Instruction:    "Convert to watercolor painting",
//		ReferenceImage: img.Base64(),
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = result.SaveTo("watercolor.png")
func (s *Service) Edit(ctx context.Context, params *EditParams) (*types.Result, error) {
	if params == nil {
		return nil, validator.ErrEmptyInstruction
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(ctx, &transport.Request{
		Method:     http.MethodPost,
		Path:       "/v1/image/edit",
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

// EditRaw modifies an image and returns raw bytes.
//
// Example:
//
//	img, _ := types.NewImageFromFile("photo.jpg")
//	result, err := client.Images.EditRaw(ctx, &image.EditParams{
//		Instruction:    "Add vintage filter",
//		ReferenceImage: img.Base64(),
//		Version:        types.VersionLatestFast,
//	}, types.FormatJPEG)
func (s *Service) EditRaw(ctx context.Context, params *EditParams, format types.OutputFormat) (*types.RawResult, error) {
	if params == nil {
		return nil, validator.ErrEmptyInstruction
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}

	if format == "" || format == types.FormatJSON {
		format = types.FormatPNG
	}

	resp, err := s.transport.DoRaw(ctx, &transport.Request{
		Method:     http.MethodPost,
		Path:       "/v1/image/edit",
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
