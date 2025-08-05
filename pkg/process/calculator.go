package process

import (
	"math"

	"github.com/Delakzz/opg-analysis/internal/pos"
)

type calculator struct {
	profitPercent   float64
	maxLossPerTrade float64
}

func (c *calculator) Calculate(gapPercent, openingPrice float64) pos.Position {
	closingPrice := openingPrice / (1 + gapPercent)
	gapValue := closingPrice - openingPrice
	profitFromGap := c.profitPercent * gapValue

	stopLoss := openingPrice - profitFromGap
	takeProfit := openingPrice + profitFromGap

	shares := int(c.maxLossPerTrade / math.Abs(stopLoss-openingPrice))
	profit := math.Abs(openingPrice-takeProfit) * float64(shares)
	profit = math.Round(profit*100) / 100 // round to 2 decimal places

	return pos.Position{
		EntryPrice:      math.Round(openingPrice*100) / 100,
		Shares:          shares,
		TakeProfitPrice: math.Round(takeProfit*100) / 100,
		StopLossPrice:   math.Round(stopLoss*100) / 100,
		Profit:          math.Round(profit*100) / 100,
	}
}

func NewCalculator(accountBalance, lossTolerance, profitPercent float64) pos.Calculator {
	maxLossPerTrade := accountBalance * lossTolerance
	return &calculator{
		profitPercent:   profitPercent,
		maxLossPerTrade: maxLossPerTrade,
	}
}
