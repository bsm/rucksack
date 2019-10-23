/*
Package runtime collects runtime metrics, if included.

	package main

	import(
		"github.com/bsm/rucksack/met"
		_ "github.com/bsm/rucksack/met/runtime"
	)

	func main() {
		sleep := make(chan struct{})
		<-sleep
	}

Run with:

	MET_NAME=myapp MET_RUNTIME=mem,heap,gc go run main.go

*/
package runtime

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bsm/instruments"
	"github.com/bsm/rucksack/v4/met"
)

var levels = make(map[string]bool)

func init() {
	for _, lev := range strings.Split(os.Getenv("MET_RUNTIME"), ",") {
		if lev != "" {
			levels[lev] = true
		}
	}
	met.Subscribe(reporter{})
}

var (
	memStats runtime.MemStats

	frees       uint64
	lookups     uint64
	mallocs     uint64
	numGC       uint32
	numCgoCalls int64
)

// Report runtime metrics
type reporter struct{}

func (reporter) Prep() error {
	if len(levels) > 0 {
		t := time.Now()
		runtime.ReadMemStats(&memStats) // This takes 50-200us.
		met.Timer("runtime.readstats", nil).Since(t)
	}
	if levels["mem"] || levels["all"] {
		reportMemStats()
	}
	if levels["heap"] || levels["all"] {
		reportHeapStats()
	}
	if levels["gc"] || levels["all"] {
		reportGCStats()
	}

	curNumCgoCalls := runtime.NumCgoCall()
	met.Gauge("runtime.cgocalls", nil).Update(float64(curNumCgoCalls - numCgoCalls))
	met.Gauge("runtime.goroutines", nil).Update(float64(runtime.NumGoroutine()))
	numCgoCalls = curNumCgoCalls
	return nil
}
func (reporter) Discrete(_ string, _ []string, _ float64) error                { return nil }
func (reporter) Sample(_ string, _ []string, _ instruments.Distribution) error { return nil }
func (reporter) Flush() error                                                  { return nil }

func reportMemStats() {
	met.Gauge("runtime.mem.alloc", nil).Update(float64(memStats.Alloc))
	met.Gauge("runtime.mem.alloc.total", nil).Update(float64(memStats.TotalAlloc))
	met.Gauge("runtime.mem.sys", nil).Update(float64(memStats.Sys))
	met.RatePerSec("runtime.mem.lookups", nil).Update(float64(memStats.Lookups - lookups))
	met.RatePerSec("runtime.mem.mallocs", nil).Update(float64(memStats.Mallocs - mallocs))
	met.RatePerSec("runtime.mem.frees", nil).Update(float64(memStats.Frees - frees))

	frees = memStats.Frees
	lookups = memStats.Lookups
	mallocs = memStats.Mallocs
}

func reportHeapStats() {
	met.Gauge("runtime.heap.alloc", nil).Update(float64(memStats.HeapAlloc))
	met.Gauge("runtime.heap.sys", nil).Update(float64(memStats.HeapSys))
	met.Gauge("runtime.heap.idle", nil).Update(float64(memStats.HeapIdle))
	met.Gauge("runtime.heap.inuse", nil).Update(float64(memStats.HeapInuse))
	met.Gauge("runtime.heap.released", nil).Update(float64(memStats.HeapReleased))
	met.Gauge("runtime.heap.objects", nil).Update(float64(memStats.HeapObjects))
}

func reportGCStats() {
	met.Gauge("runtime.gc.next", nil).Update(float64(memStats.NextGC))
	met.Gauge("runtime.gc.last", nil).Update(float64(memStats.LastGC))
	met.RatePerMin("runtime.gc.num", nil).Update(float64(memStats.NumGC - numGC))
	met.Gauge("runtime.gc.cpu", nil).Update(memStats.GCCPUFraction * 1000)

	i := numGC % uint32(len(memStats.PauseNs))
	pauseGC := met.Timer("runtime.gc.pause", nil)
	if memStats.NumGC-numGC >= uint32(len(memStats.PauseNs)) {
		for i = 0; i < uint32(len(memStats.PauseNs)); i++ {
			pauseGC.Update(time.Duration(memStats.PauseNs[i]))
		}
	} else {
		ii := memStats.NumGC % uint32(len(memStats.PauseNs))
		if i > ii {
			for ; i < uint32(len(memStats.PauseNs)); i++ {
				pauseGC.Update(time.Duration(memStats.PauseNs[i]))
			}
			i = 0
		}
		for ; i < ii; i++ {
			pauseGC.Update(time.Duration(memStats.PauseNs[i]))
		}
	}
	met.Gauge("runtime.gc.pause.total", nil).Update(float64(memStats.PauseTotalNs) / 1e6)

	numGC = memStats.NumGC
}
