package trade

import (
	"github.com/Delakzz/opg-analysis/internal/news"
	"github.com/Delakzz/opg-analysis/internal/pos"
)

type Selection struct {
	Ticker string
	pos.Position
	Articles []news.Article
}

type Deliverer interface {
	Deliver(selections []Selection) error
}
