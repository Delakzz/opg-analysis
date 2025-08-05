package process

import (
	"math"
	"slices"

	"github.com/Delakzz/opg-analysis/internal/raw"
)

type filterer struct {
	minGap float64
}

func (n *filterer) Filter(candidates []raw.Stock) (filtered []raw.Stock) {
	filtered = slices.DeleteFunc(candidates, func(stock raw.Stock) bool {
		return math.Abs(stock.Gap) < n.minGap
	})
	return
}

func NewFilterer(minGap float64) raw.Filterer {
	return &filterer{
		minGap: minGap,
	}
}
