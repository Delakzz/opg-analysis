package news

import "time"

type Article struct {
	PublishOn time.Time
	Headline  string
}

type Fetcher interface {
	Fetch(ticker string) ([]Article, error)
}
