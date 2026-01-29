package image

import (
	"context"
	"sync"

	"github.com/shamspias/reve-go/types"
)

// BatchConfig configures batch operations.
type BatchConfig struct {
	// Concurrency is the max concurrent requests.
	// Default: 5
	Concurrency int

	// StopOnError stops on first error.
	// Default: false
	StopOnError bool
}

// DefaultBatchConfig returns default configuration.
func DefaultBatchConfig() *BatchConfig {
	return &BatchConfig{
		Concurrency: 5,
		StopOnError: false,
	}
}

// BatchResult represents a single batch result.
type BatchResult struct {
	// Index is the position in the request slice.
	Index int

	// Result is the successful result.
	Result *types.Result

	// Error is the error if failed.
	Error error
}

// BatchCreate executes multiple create requests concurrently.
//
// Example:
//
//	requests := []*image.CreateParams{
//		{Prompt: "A red apple"},
//		{Prompt: "A green pear"},
//		{Prompt: "Purple grapes"},
//	}
//
//	results := client.Images.BatchCreate(ctx, requests, &image.BatchConfig{
//		Concurrency: 3,
//	})
//
//	for _, r := range results {
//		if r.Error != nil {
//			log.Printf("Request %d failed: %v", r.Index, r.Error)
//		} else {
//			r.Result.SaveTo(fmt.Sprintf("fruit_%d.png", r.Index))
//		}
//	}
func (s *Service) BatchCreate(ctx context.Context, params []*CreateParams, config *BatchConfig) []BatchResult {
	if config == nil {
		config = DefaultBatchConfig()
	}

	results := make([]BatchResult, len(params))
	var wg sync.WaitGroup
	sem := make(chan struct{}, config.Concurrency)
	var stopFlag bool
	var stopMu sync.Mutex

	for i, p := range params {
		if ctx.Err() != nil {
			results[i] = BatchResult{Index: i, Error: ctx.Err()}
			continue
		}

		stopMu.Lock()
		if stopFlag {
			stopMu.Unlock()
			results[i] = BatchResult{Index: i, Error: context.Canceled}
			continue
		}
		stopMu.Unlock()

		wg.Add(1)
		go func(idx int, req *CreateParams) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				results[idx] = BatchResult{Index: idx, Error: ctx.Err()}
				return
			}

			result, err := s.Create(ctx, req)
			results[idx] = BatchResult{Index: idx, Result: result, Error: err}

			if err != nil && config.StopOnError {
				stopMu.Lock()
				stopFlag = true
				stopMu.Unlock()
			}
		}(i, p)
	}

	wg.Wait()
	return results
}

// BatchEdit executes multiple edit requests concurrently.
//
// Example:
//
//	img, _ := types.NewImageFromFile("photo.jpg")
//	requests := []*image.EditParams{
//		{Instruction: "Make warmer", ReferenceImage: img.Base64()},
//		{Instruction: "Make cooler", ReferenceImage: img.Base64()},
//		{Instruction: "Add vintage", ReferenceImage: img.Base64()},
//	}
//
//	results := client.Images.BatchEdit(ctx, requests, nil)
func (s *Service) BatchEdit(ctx context.Context, params []*EditParams, config *BatchConfig) []BatchResult {
	if config == nil {
		config = DefaultBatchConfig()
	}

	results := make([]BatchResult, len(params))
	var wg sync.WaitGroup
	sem := make(chan struct{}, config.Concurrency)
	var stopFlag bool
	var stopMu sync.Mutex

	for i, p := range params {
		if ctx.Err() != nil {
			results[i] = BatchResult{Index: i, Error: ctx.Err()}
			continue
		}

		stopMu.Lock()
		if stopFlag {
			stopMu.Unlock()
			results[i] = BatchResult{Index: i, Error: context.Canceled}
			continue
		}
		stopMu.Unlock()

		wg.Add(1)
		go func(idx int, req *EditParams) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				results[idx] = BatchResult{Index: idx, Error: ctx.Err()}
				return
			}

			result, err := s.Edit(ctx, req)
			results[idx] = BatchResult{Index: idx, Result: result, Error: err}

			if err != nil && config.StopOnError {
				stopMu.Lock()
				stopFlag = true
				stopMu.Unlock()
			}
		}(i, p)
	}

	wg.Wait()
	return results
}

// BatchRemix executes multiple remix requests concurrently.
func (s *Service) BatchRemix(ctx context.Context, params []*RemixParams, config *BatchConfig) []BatchResult {
	if config == nil {
		config = DefaultBatchConfig()
	}

	results := make([]BatchResult, len(params))
	var wg sync.WaitGroup
	sem := make(chan struct{}, config.Concurrency)
	var stopFlag bool
	var stopMu sync.Mutex

	for i, p := range params {
		if ctx.Err() != nil {
			results[i] = BatchResult{Index: i, Error: ctx.Err()}
			continue
		}

		stopMu.Lock()
		if stopFlag {
			stopMu.Unlock()
			results[i] = BatchResult{Index: i, Error: context.Canceled}
			continue
		}
		stopMu.Unlock()

		wg.Add(1)
		go func(idx int, req *RemixParams) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				results[idx] = BatchResult{Index: idx, Error: ctx.Err()}
				return
			}

			result, err := s.Remix(ctx, req)
			results[idx] = BatchResult{Index: idx, Result: result, Error: err}

			if err != nil && config.StopOnError {
				stopMu.Lock()
				stopFlag = true
				stopMu.Unlock()
			}
		}(i, p)
	}

	wg.Wait()
	return results
}

// SuccessCount returns number of successful results.
func SuccessCount(results []BatchResult) int {
	count := 0
	for _, r := range results {
		if r.Error == nil {
			count++
		}
	}
	return count
}

// ErrorCount returns number of failed results.
func ErrorCount(results []BatchResult) int {
	count := 0
	for _, r := range results {
		if r.Error != nil {
			count++
		}
	}
	return count
}

// Successful returns only successful results.
func Successful(results []BatchResult) []*types.Result {
	var out []*types.Result
	for _, r := range results {
		if r.Error == nil {
			out = append(out, r.Result)
		}
	}
	return out
}

// Errors returns all errors.
func Errors(results []BatchResult) []error {
	var errs []error
	for _, r := range results {
		if r.Error != nil {
			errs = append(errs, r.Error)
		}
	}
	return errs
}
