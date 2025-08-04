package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
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
}

func main() {

	stocks, err := Load("./opg.csv")
	if err != nil {
		fmt.Print(err)
		return
	}

	slices.DeleteFunc(stocks, func(s Stock) bool {
		return math.Abs(s.Gap) < 0.1
	})

	var selections []Selection

	for _, stock := range stocks {
		position := Calculate(stock.Gap, stock.OpeningPrice)

		sel := Selection{
			Ticker:   stock.Ticker,
			Position: position,
		}

		selections = append(selections, sel)
	}
}
