/*
Package jaeger reports metrics to Jaeger, if included.

  package main

  import(
    _ "github.com/bsm/rucksack/trace/jaeger"
  )

  func main() {
    sleep := make(chan struct{})
    <-sleep
  }

Run with:

  TRACE_NAME=myapp TRACE_JAEGER=10.0.0.1 go run main.go

*/
package jaeger

import (
	"io"
	"runtime"

	"github.com/bsm/rucksack"
	"github.com/bsm/rucksack/log"
	"github.com/uber/jaeger-client-go/config"
)

func init() {
	if name, host := rucksack.Env("TRACE_NAME", "APP_NAME"), rucksack.Env("TRACE_JAEGER"); name != "" && host != "" {
		cfg := config.Configuration{
			Sampler: &config.SamplerConfig{
				SamplingServerURL: "http://" + host + ":5778/sampling",
			},
			Reporter: &config.ReporterConfig{
				LocalAgentHostPort: host + ":5775",
			},
		}

		closer, err := cfg.InitGlobalTracer(name)
		if err != nil {
			log.Errorf("trace/jaeger: could not initialize tracer: %v", err)
			return
		}
		runtime.SetFinalizer(closer, func(c io.Closer) { _ = c.Close() })
	}
}
