package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"
)

type Stock struct {
	Ticker       string
	Gap          float64
	OpeningPrice float64
}

func Load(path string) ([]Stock, error) {
	f, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rows = slices.Delete(rows, 0, 1)

	var stocks []Stock

	for _, row := range rows {
		ticker := row[0]
		gap, err := strconv.ParseFloat(row[1], 64)

		if err != nil {
			continue
		}

		openingPrice, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}

		stocks = append(stocks, Stock{
			Ticker:       ticker,
			Gap:          gap,
			OpeningPrice: openingPrice,
		})
	}

	return stocks, nil
}

var (
	accountBalance  = 100000.0                       // how much money in trading account
	lossTolerance   = .02                            // what percentage of that balance I can tolerate losing
	maxLossPerTrade = accountBalance * lossTolerance // maximum maount i can tolerate losing
	profitPercent   = .8                             // percentage of the gap I want to take as profit
)

type Position struct {
	EntryPrice      float64 // the price at which to buy or sell
	Shares          int     // how many shares to buy or sell
	TakeProfitPrice float64 // the price at which to exit and take my profit
	StopLossPrice   float64 // the price at which to stop my loss if the stock doesn't go our way
	Profit          float64 // expected final profit
}

func Calculate(gapPercent, openingPrice float64) Position {
	closingPrice := openingPrice / (1 + gapPercent)
	gapValue := closingPrice - openingPrice
	profitFromGap := profitPercent * gapValue

	stopLoss := openingPrice - profitFromGap
	takeProfit := openingPrice + profitFromGap

	shares := int(maxLossPerTrade / math.Abs(stopLoss-openingPrice))

	profit := math.Abs(openingPrice-takeProfit) * float64(shares)
	profit = math.Round(profit*100) / 100 // round to 2 decimal places

	return Position{
		EntryPrice:      math.Round(openingPrice*100) / 100,
		Shares:          shares,
		TakeProfitPrice: math.Round(takeProfit*100) / 100,
		StopLossPrice:   math.Round(stopLoss*100) / 100,
		Profit:          math.Round(profit*100) / 100,
	}
}

type Selection struct {
	Ticker string
	Position
	Articles []Article
}

const (
	url          = "https://seeking-alpha.p.rapidapi.com/news/v2/list-by-symbol?id=AAPL&size=5&id="
	apiKeyHeader = "x-rapidapi-key"
	apiKey       = "bfe66c058dmsha2f18e51b86c7d1p1a0fbejsnd6258bc8b219"
)

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

type Article struct {
	PublishOn time.Time
	Headline  string
}

func FetchNews(ticker string) ([]Article, error) {
	req, err := http.NewRequest(http.MethodGet, url+ticker, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add(apiKeyHeader, apiKey)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("error fetching news for %s: %s", ticker, resp.Status)
	}

	res := &SeekingAlphaResponse{}
	json.NewDecoder(resp.Body).Decode(res)

	var articles []Article

	for _, item := range res.Data {
		art := Article{
			PublishOn: item.Attributes.PublishOn,
			Headline:  item.Attributes.Title,
		}
		articles = append(articles, art)
	}

	return articles, nil
}

func Deliver(filePath string, selections []Selection) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(selections)
	if err != nil {
		return fmt.Errorf("error encoding selections to JSON: %w", err)
	}

	return nil
}

func main() {

	stocks, err := Load("./opg.csv")
	if err != nil {
		fmt.Print(err)
		return
	}

	slices.DeleteFunc(stocks, func(s Stock) bool {
		return math.Abs(s.Gap) < .1
	})

	var selections []Selection

	for _, stock := range stocks {
		position := Calculate(stock.Gap, stock.OpeningPrice)

		articles, err := FetchNews(stock.Ticker)
		if err != nil {
			fmt.Printf("Error fetching news for %s: %v\n", stock.Ticker, err)
			continue
		} else {
			log.Printf("Found %d articles for %s\n", len(articles), stock.Ticker)
		}
		sel := Selection{
			Ticker:   stock.Ticker,
			Position: position,
			Articles: articles,
		}

		selections = append(selections, sel)
	}

	outputPath := "./opg.json"

	err = Deliver(outputPath, selections)
	if err != nil {
		log.Printf("Error writing output %s: %v\n", outputPath, err)
		return
	}

	log.Printf("Wrote %d selections to %s\n", len(selections), outputPath)
}
