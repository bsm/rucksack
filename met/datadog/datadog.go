package datadog

import (
	"os"

	"github.com/bsm/instruments/datadog"
	"github.com/bsm/rucksack/met"
)

func init() {
	if token := os.Getenv("MET_DATADOG"); token != "" {
		met.Subscribe(datadog.New(token))
	}
}
