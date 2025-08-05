package cmd

import (
	"fmt"
	"log"

	"github.com/Delakzz/opg-analysis/internal/news"
	"github.com/Delakzz/opg-analysis/internal/pos"
	"github.com/Delakzz/opg-analysis/internal/raw"
	"github.com/Delakzz/opg-analysis/internal/trade"
)

func Run(ldr raw.Loader, f raw.Filterer, c pos.Calculator, fet news.Fetcher, del trade.Deliverer) error {
	stocks, err := ldr.Load()
	if err != nil {
		return fmt.Errorf("failed to load stocks: %w", err)
	}

	stocks = f.Filter(stocks)

	selectionsChannel := make(chan trade.Selection, len(stocks))

	for _, stock := range stocks {
		go func(s raw.Stock, selected chan<- trade.Selection) {
			position := c.Calculate(s.Gap, s.OpeningPrice)

			articles, err := fet.Fetch(s.Ticker)
			if err != nil {
				log.Printf("error loading news about %s, %v", s.Ticker, err)
				selected <- trade.Selection{}
				return
			} else {
				log.Printf("found %d articles about %s", len(articles), s.Ticker)
			}

			sel := trade.Selection{
				Ticker:   s.Ticker,
				Position: position,
				Articles: articles,
			}

			selected <- sel
		}(stock, selectionsChannel)
	}

	var selections []trade.Selection

	for sel := range selectionsChannel {
		selections = append(selections, sel)
		if len(selections) == len(stocks) {
			close(selectionsChannel)
		}
	}

	err = del.Deliver(selections)
	if err != nil {
		return fmt.Errorf("error delivering selections: %w", err)
	}

	return nil
}
