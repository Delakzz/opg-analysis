package json

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Delakzz/opg-analysis/internal/trade"
)

type deliverer struct {
	filepath string
}

func (d *deliverer) Deliver(selections []trade.Selection) error {
	file, err := os.Create(d.filepath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", d.filepath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(selections)
	if err != nil {
		return fmt.Errorf("error encoding selections to JSON: %w", err)
	}

	log.Printf("Selections successfully delivered to %s", d.filepath)
	return nil
}

func NewDeliverer(filepath string) trade.Deliverer {
	return &deliverer{filepath: filepath}
}
