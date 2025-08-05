package salpha

import "time"

type attributes struct {
	PublishOn time.Time `json:"publishOn"`
	Title     string    `json:"title"`
}
type seekingAlphaNews struct {
	Attributes attributes `json:"attributes"`
}
type SeekingAlphaResponse struct {
	Data []seekingAlphaNews `json:"data"`
}
