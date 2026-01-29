package reve

// Re-export commonly used types for convenience.
// Users can import just "github.com/shamspias/reve-go" and use these directly.

import (
	"github.com/shamspias/reve-go/image"
	"github.com/shamspias/reve-go/types"
)

// Type aliases for convenience.
type (
	// AspectRatio represents image aspect ratios.
	AspectRatio = types.AspectRatio

	// ModelVersion represents model versions.
	ModelVersion = types.ModelVersion

	// OutputFormat represents response formats.
	OutputFormat = types.OutputFormat

	// Postprocess represents postprocessing operations.
	Postprocess = types.Postprocess

	// Image represents an image for API operations.
	Image = types.Image

	// Result represents an image generation result.
	Result = types.Result

	// RawResult represents a raw binary result.
	RawResult = types.RawResult

	// CreateParams is parameters for image creation.
	CreateParams = image.CreateParams

	// EditParams is parameters for image editing.
	EditParams = image.EditParams

	// RemixParams is parameters for image remixing.
	RemixParams = image.RemixParams

	// BatchConfig configures batch operations.
	BatchConfig = image.BatchConfig

	// BatchResult represents a batch operation result.
	BatchResult = image.BatchResult

	// Cost represents an estimated cost.
	Cost = image.Cost
)

// Aspect ratio constants.
const (
	Ratio16x9 = types.Ratio16x9
	Ratio9x16 = types.Ratio9x16
	Ratio3x2  = types.Ratio3x2
	Ratio2x3  = types.Ratio2x3
	Ratio4x3  = types.Ratio4x3
	Ratio3x4  = types.Ratio3x4
	Ratio1x1  = types.Ratio1x1
	RatioAuto = types.RatioAuto
)

// Model version constants.
const (
	VersionLatest            = types.VersionLatest
	VersionLatestFast        = types.VersionLatestFast
	VersionCreate20250915    = types.VersionCreate20250915
	VersionEdit20250915      = types.VersionEdit20250915
	VersionEditFast20251030  = types.VersionEditFast20251030
	VersionRemix20250915     = types.VersionRemix20250915
	VersionRemixFast20251030 = types.VersionRemixFast20251030
)

// Output format constants.
const (
	FormatJSON = types.FormatJSON
	FormatPNG  = types.FormatPNG
	FormatJPEG = types.FormatJPEG
	FormatWebP = types.FormatWebP
)

// Helper functions re-exported for convenience.
var (
	// NewImage creates an Image from bytes.
	NewImage = types.NewImage

	// NewImageFromBase64 creates an Image from base64.
	NewImageFromBase64 = types.NewImageFromBase64

	// NewImageFromFile loads an Image from file.
	NewImageFromFile = types.NewImageFromFile

	// Ref creates an image reference tag.
	Ref = types.Ref

	// Upscale creates an upscale operation.
	Upscale = types.Upscale

	// RemoveBackground creates a background removal operation.
	RemoveBackground = types.RemoveBackground

	// DetectFormat detects format from file path.
	DetectFormat = types.DetectFormat

	// EstimateCreate estimates create cost.
	EstimateCreate = image.EstimateCreate

	// EstimateEdit estimates edit cost.
	EstimateEdit = image.EstimateEdit

	// EstimateRemix estimates remix cost.
	EstimateRemix = image.EstimateRemix

	// FormatCredits formats credits as string.
	FormatCredits = image.FormatCredits

	// DefaultBatchConfig returns default batch config.
	DefaultBatchConfig = image.DefaultBatchConfig

	// SuccessCount returns successful results count.
	SuccessCount = image.SuccessCount

	// ErrorCount returns failed results count.
	ErrorCount = image.ErrorCount

	// Successful returns successful results.
	Successful = image.Successful

	// Errors returns all errors from batch.
	Errors = image.Errors
)
