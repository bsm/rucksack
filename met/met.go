// Package met provides a 12-factor convenience wrapper around instruments
package met

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bsm/instruments"
)

// NewRegistry returns a new, custom registry
func NewRegistry(name string) *instruments.Registry {
	tags := make([]string, len(defaultTags))
	copy(tags, defaultTags)
	return instruments.New(time.Minute, name+".", tags...)
}

// Hostname returns the parsed hostname
func Hostname() string { return hostname }

// --------------------------------------------------------------------

// Subscribe attaches reporters/hooks to the met registry
func Subscribe(rep instruments.Reporter) {
	registry.Subscribe(rep)
}

// AddTags add tags to default registry
func AddTags(tags ...string) {
	registry.AddTags(tags...)
}

// Convenience accessors to default registry metrics

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
func Reservoir(name string, tags []string, size int) *instruments.Reservoir {
	return registry.Reservoir(name, tags, size)
}
func Timer(name string, tags []string, size int) *instruments.Timer {
	return registry.Timer(name, tags, size)
}

// --------------------------------------------------------------------

var (
	registry    = instruments.NewUnstarted("")
	hostname    string
	defaultTags []string
)

func init() {
	// Parse the name of the app to meter
	name := os.Getenv("MET_NAME")
	if name == "" {
		return
	}

	// Parse hostname
	hostname = os.Getenv("HOST")
	if hostname == "" {
		hostname, _ = os.Hostname()
	}
	if pos := strings.Index(hostname, "."); pos > -1 {
		hostname = hostname[:pos]
	}

	if hostname != "" {
		defaultTags = append(defaultTags, "host:"+hostname)
	}
	if port := os.Getenv("PORT"); port != "" {
		defaultTags = append(defaultTags, "port:"+port)
	}
	if othr := os.Getenv("MET_TAGS"); othr != "" {
		defaultTags = append(defaultTags, strings.Split(othr, ",")...)
	}

	// Create registry
	registry = NewRegistry(name)
	runtime.SetFinalizer(registry, func(r *instruments.Registry) { _ = r.Close() })
}
