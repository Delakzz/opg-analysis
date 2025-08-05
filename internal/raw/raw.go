package raw

type Stock struct {
	Ticker       string
	Gap          float64
	OpeningPrice float64
}

type Loader interface {
	Load() ([]Stock, error)
}

type Filterer interface {
	Filter([]Stock) []Stock
}
