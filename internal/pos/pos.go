package pos

type Position struct {
	EntryPrice      float64
	Shares          int
	TakeProfitPrice float64
	StopLossPrice   float64
	Profit          float64
}

type Calculator interface {
	Calculate(gapPercent, openingPrice float64) Position
}
