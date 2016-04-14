// Package met provides a 12-factor convenience wrapper around instruments
package met

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bsm/instruments"
)

var registry = instruments.NewUnstarted("")

// Subscribe attaches reporters/hooks to the met registry
func Subscribe(rep instruments.Reporter) {
	registry.Subscribe(rep)
}

// Convenience accessors to metrics

func Counter(name string, tags []string) *instruments.Counter {
	return registry.Counter(name, tags)
}
func Gauge(name string, tags []string) *instruments.Gauge {
	return registry.Gauge(name, tags)
}
func RatePerSec(name string, tags []string) *instruments.Rate {
	return registry.RateScale(name, tags, time.Second)
}
func RatePerMin(name string, tags []string) *instruments.Rate {
	return registry.RateScale(name, tags, time.Minute)
}
func RateScale(name string, tags []string, d time.Duration) *instruments.Rate {
	return registry.RateScale(name, tags, d)
}
func Timer(name string, tags []string, size int64) *instruments.Timer {
	return registry.Timer(name, tags, size)
}

// --------------------------------------------------------------------

func init() {
	// Parse the name of the app to meter
	name := os.Getenv("MET_NAME")
	if name == "" {
		return
	}

	// Parse tags
	var tags []string
	if host := hostname(); host != "" {
		tags = append(tags, "host:"+host)
	}
	if port := os.Getenv("PORT"); port != "" {
		tags = append(tags, "port:"+port)
	}
	if othr := os.Getenv("MET_TAGS"); othr != "" {
		tags = append(tags, strings.Split(othr, ",")...)
	}

	// Create registry
	registry = instruments.New(time.Minute, name+".", tags...)
	runtime.SetFinalizer(registry, func(r *instruments.Registry) { _ = r.Close() })
}

func hostname() string {
	host := os.Getenv("HOST")
	if host == "" {
		host, _ = os.Hostname()
	}
	if pos := strings.Index(host, "."); pos > -1 {
		host = host[:pos]
	}
	return host
}
