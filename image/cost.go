package image

import (
	"fmt"

	"github.com/shamspias/reve-go/types"
)

// Credit costs.
const (
	CostCreate    = 18
	CostEdit      = 30
	CostEditFast  = 5
	CostRemix     = 30
	CostRemixFast = 5
	CostPerCredit = 0.04 / 30.0 // ~$0.00133
)

// Cost represents an estimated cost.
type Cost struct {
	BaseCredits  int
	TotalCredits int
	EstimatedUSD float64
}

// String returns a readable string.
func (c Cost) String() string {
	return fmt.Sprintf("%d credits (~$%.4f)", c.TotalCredits, c.EstimatedUSD)
}

// EstimateCreate estimates create cost.
//
// Example:
//
//	cost := image.EstimateCreate(1, nil)
//	fmt.Println(cost) // "18 credits (~$0.0240)"
func EstimateCreate(scaling float64, postprocess []types.Postprocess) Cost {
	return estimate(CostCreate, scaling, postprocess)
}

// EstimateEdit estimates edit cost.
//
// Example:
//
//	cost := image.EstimateEdit(false, 1, nil)  // Standard
//	fmt.Println(cost) // "30 credits (~$0.0400)"
//
//	cost = image.EstimateEdit(true, 1, nil)   // Fast
//	fmt.Println(cost) // "5 credits (~$0.0067)"
func EstimateEdit(fast bool, scaling float64, postprocess []types.Postprocess) Cost {
	base := CostEdit
	if fast {
		base = CostEditFast
	}
	return estimate(base, scaling, postprocess)
}

// EstimateRemix estimates remix cost.
func EstimateRemix(fast bool, scaling float64, postprocess []types.Postprocess) Cost {
	base := CostRemix
	if fast {
		base = CostRemixFast
	}
	return estimate(base, scaling, postprocess)
}

func estimate(base int, scaling float64, postprocess []types.Postprocess) Cost {
	total := base

	if scaling > 1 {
		total = int(float64(total) * scaling)
	}

	for _, pp := range postprocess {
		switch pp.Process {
		case types.ProcessUpscale:
			total += pp.UpscaleFactor * 2
		case types.ProcessRemoveBackground:
			total += 2
		}
	}

	return Cost{
		BaseCredits:  base,
		TotalCredits: total,
		EstimatedUSD: float64(total) * CostPerCredit,
	}
}

// FormatCredits formats credits as readable string.
func FormatCredits(credits int) string {
	return fmt.Sprintf("%d credits (~$%.4f)", credits, float64(credits)*CostPerCredit)
}
