package salpha

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Delakzz/opg-analysis/internal/news"
)

const (
	urlPath      = "/news/v2/list-by-symbol"
	apiKeyHeader = "bfe66c058dmsha2f18e51b86c7d1p1a0fbejsnd6258bc8b219"
	pageSize     = 5
)

type client struct {
	baseURL string
	apiKey  string
}

func (c *client) Fetch(ticker string) ([]news.Article, error) {
	url, err := c.buildURL(ticker)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add(apiKeyHeader, c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("unsuccessful status code %d received", resp.StatusCode)
	}

	return c.parse(resp)
}

func (c *client) parse(resp *http.Response) ([]news.Article, error) {
	res := &SeekingAlphaResponse{}
	err := json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, err
	}

	var articles []news.Article
	for _, item := range res.Data {
		art := news.Article{
			PublishOn: item.Attributes.PublishOn,
			Headline:  item.Attributes.Title,
		}
		articles = append(articles, art)
	}

	return articles, nil
}

func (c *client) buildURL(ticker string) (string, error) {
	parsedURL, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}

	parsedURL.Path += urlPath

	params := url.Values{}
	params.Add("size", fmt.Sprint(pageSize))
	params.Add("id", ticker)
	parsedURL.RawQuery = params.Encode()

	return parsedURL.String(), nil
}

func NewClient(baseURL, apiKey string) news.Fetcher {
	return &client{baseURL: baseURL, apiKey: apiKey}
}
