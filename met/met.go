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
var hostname string

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

// Hostname returns the parsed hostname
func Hostname() string { return hostname }

// AddTags alows to add global tags
func AddTags(tags ...string) { registry.AddTags(tags...) }

// --------------------------------------------------------------------

func init() {
	// Parse the name of the app to meter
	name := os.Getenv("MET_NAME")
	if name == "" {
		return
	}

	// Parse tags
	var tags []string

	// Parse hostname
	hostname = os.Getenv("HOST")
	if hostname == "" {
		hostname, _ = os.Hostname()
	}
	if pos := strings.Index(hostname, "."); pos > -1 {
		hostname = hostname[:pos]
	}

	if hostname != "" {
		tags = append(tags, "host:"+hostname)
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
