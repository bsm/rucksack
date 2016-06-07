/*
Package datadog reports metrics to datadog, if included.

  package main

  import(
    "github.com/bsm/rucksack/met"
    _ "github.com/bsm/rucksack/met/datadog"
  )

  func main() {
    sleep := make(chan struct{})
    <-sleep
  }

Run with:

  MET_NAME=myapp MET_DATADOG=mytoken go run main.go

*/
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
