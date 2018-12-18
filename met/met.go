// Package met provides a 12-factor convenience wrapper around instruments
package met

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bsm/instruments"
	"github.com/bsm/rucksack"
)

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
func Reservoir(name string, tags []string) *instruments.Reservoir {
	return registry.Reservoir(name, tags)
}
func Timer(name string, tags []string) *instruments.Timer {
	return registry.Timer(name, tags)
}

// --------------------------------------------------------------------

var (
	registry = instruments.NewUnstarted("")
	hostname string
)

func init() {
	// Parse the name of the app to meter
	name := rucksack.Env("MET_NAME", "APP_NAME")
	if name == "" {
		return
	}

	// Parse hostname
	hostname = rucksack.Env("HOST")
	if hostname == "" {
		hostname, _ = os.Hostname()
	}
	if pos := strings.Index(hostname, "."); pos > -1 {
		hostname = hostname[:pos]
	}

	// Parse tags
	var tags []string
	if hostname != "" {
		tags = append(tags, "host:"+hostname)
	}
	if port := rucksack.Env("PORT"); port != "" {
		tags = append(tags, "port:"+port)
	}
	tags = append(tags, rucksack.Tags(rucksack.Env("MET_TAGS", "APP_TAGS"))...)

	// Parse interval
	interval, _ := time.ParseDuration(rucksack.Env("MET_REPORT_INTERVAL"))
	if interval == 0 {
		interval = 30 * time.Second
	}

	// Create registry
	registry = instruments.New(interval, name+".", tags...)
	runtime.SetFinalizer(registry, func(r *instruments.Registry) { _ = r.Close() })
}
