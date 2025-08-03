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

func main() {

	stocks, err := Load("./opg.csv")
	if err != nil {
		fmt.Print(err)
		return
	}

	slices.DeleteFunc(stocks, func(s Stock) bool {
		return math.Abs(s.Gap) < 0.1
	})

}
