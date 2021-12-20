package agent

import (
	"log"
	"math/rand"
	"runtime"
)

type Gauge float64

type Counter int64

type Metrics struct {
	// Alloc is bytes of allocated heap objects.
	Alloc Gauge

	// BuckHashSys is bytes of memory in profiling bucket hash tables.
	BuckHashSys Gauge

	// Frees is the cumulative count of heap objects freed
	// in this size class.
	Frees Gauge

	// GCCPUFraction is the fraction of this program's available
	// CPU time used by the GC since the program started.
	//
	// GCCPUFraction is expressed as a number between 0 and 1,
	// where 0 means GC has consumed none of this program's CPU. A
	// program's available CPU time is defined as the integral of
	// GOMAXPROCS since the program started. That is, if
	// GOMAXPROCS is 2 and a program has been running for 10
	// seconds, its "available CPU" is 20 seconds. GCCPUFraction
	// does not include CPU time used for write barrier activity.
	//
	// This is the same as the fraction of CPU reported by
	// GODEBUG=gctrace=1.
	GCCPUFraction Gauge

	// GCSys is bytes of memory in garbage collection metadata.
	GCSys Gauge

	// HeapAlloc is bytes of allocated heap objects.
	//
	// "Allocated" heap objects include all reachable objects, as
	// well as unreachable objects that the garbage collector has
	// not yet freed. Specifically, HeapAlloc increases as heap
	// objects are allocated and decreases as the heap is swept
	// and unreachable objects are freed. Sweeping occurs
	// incrementally between GC cycles, so these two processes
	// occur simultaneously, and as a result HeapAlloc tends to
	// change smoothly (in contrast with the sawtooth that is
	// typical of stop-the-world garbage collectors).
	HeapAlloc Gauge

	// HeapIdle is bytes in idle (unused) spans.
	//
	// Idle spans have no objects in them. These spans could be
	// (and may already have been) returned to the OS, or they can
	// be reused for heap allocations, or they can be reused as
	// stack memory.
	//
	// HeapIdle minus HeapReleased estimates the amount of memory
	// that could be returned to the OS, but is being retained by
	// the runtime so it can grow the heap without requesting more
	// memory from the OS. If this difference is significantly
	// larger than the heap size, it indicates there was a recent
	// transient spike in live heap size.
	HeapIdle Gauge

	// HeapInuse is bytes in in-use spans.
	//
	// In-use spans have at least one object in them. These spans
	// can only be used for other objects of roughly the same
	// size.
	//
	// HeapInuse minus HeapAlloc estimates the amount of memory
	// that has been dedicated to particular size classes, but is
	// not currently being used. This is an upper bound on
	// fragmentation, but in general this memory can be reused
	// efficiently.
	HeapInuse Gauge

	// HeapObjects is the number of allocated heap objects.
	//
	// Like HeapAlloc, this increases as objects are allocated and
	// decreases as the heap is swept and unreachable objects are
	// freed.
	HeapObjects Gauge

	// HeapReleased is bytes of physical memory returned to the OS.
	//
	// This counts heap memory from idle spans that was returned
	// to the OS and has not yet been reacquired for the heap.
	HeapReleased Gauge

	// HeapSys is bytes of heap memory obtained from the OS.
	//
	// HeapSys measures the amount of virtual address space
	// reserved for the heap. This includes virtual address space
	// that has been reserved but not yet used, which consumes no
	// physical memory, but tends to be small, as well as virtual
	// address space for which the physical memory has been
	// returned to the OS after it became unused (see HeapReleased
	// for a measure of the latter).
	//
	// HeapSys estimates the largest size the heap has had.
	HeapSys Gauge

	// LastGC is the time the last garbage collection finished, as
	// nanoseconds since 1970 (the UNIX epoch).
	LastGC Gauge

	// Lookups is the number of pointer lookups performed by the
	// runtime.
	//
	// This is primarily useful for debugging runtime internals.
	Lookups Gauge

	// MCacheInuse is bytes of allocated mcache structures.
	MCacheInuse Gauge

	// MCacheSys is bytes of memory obtained from the OS for
	// mcache structures.
	MCacheSys Gauge

	// MSpanInuse is bytes of allocated mspan structures.
	MSpanInuse Gauge

	// MSpanSys is bytes of memory obtained from the OS for mspan
	// structures.
	MSpanSys Gauge

	// Mallocs is the cumulative count of heap objects allocated.
	// The number of live objects is Mallocs - Frees.
	Mallocs Gauge

	// NextGC is the target heap size of the next GC cycle.
	//
	// The garbage collector's goal is to keep HeapAlloc ≤ NextGC.
	// At the end of each GC cycle, the target for the next cycle
	// is computed based on the amount of reachable data and the
	// value of GOGC.
	NextGC Gauge

	// NumForcedGC is the number of GC cycles that were forced by
	// the application calling the GC function.
	NumForcedGC Gauge

	// NumGC is the number of completed GC cycles.
	NumGC Gauge

	// OtherSys is bytes of memory in miscellaneous off-heap
	// runtime allocations.
	OtherSys Gauge

	// PauseTotalNs is the cumulative nanoseconds in GC
	// stop-the-world pauses since the program started.
	//
	// During a stop-the-world pause, all goroutines are paused
	// and only the garbage collector can run.
	PauseTotalNs Gauge

	// StackInuse is bytes in stack spans.
	//
	// In-use stack spans have at least one stack in them. These
	// spans can only be used for other stacks of the same size.
	//
	// There is no StackIdle because unused stack spans are
	// returned to the heap (and hence counted toward HeapIdle).
	StackInuse Gauge

	// StackSys is bytes of stack memory obtained from the OS.
	//
	// StackSys is StackInuse, plus any memory obtained directly
	// from the OS for OS thread stacks (which should be minimal).
	StackSys Gauge

	// Sys is the total bytes of memory obtained from the OS.
	//
	// Sys is the sum of the XSys fields below. Sys measures the
	// virtual address space reserved by the Go runtime for the
	// heap, stacks, and other internal data structures. It's
	// likely that not all of the virtual address space is backed
	// by physical memory at any given moment, though in general
	// it all was at some point.
	Sys Gauge

	// PollCount is the number of previous polls.
	PollCount Counter

	// RandomValue is just a random value.
	RandomValue Gauge
}

func (m *Metrics) feelMemStats() {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	m.Alloc = Gauge(memStats.Alloc)
	m.BuckHashSys = Gauge(memStats.BuckHashSys)
	m.Frees = Gauge(memStats.Frees)
	m.GCCPUFraction = Gauge(memStats.GCCPUFraction)
	m.GCSys = Gauge(memStats.GCSys)
	m.HeapAlloc = Gauge(memStats.HeapAlloc)
	m.HeapIdle = Gauge(memStats.HeapIdle)
	m.HeapInuse = Gauge(memStats.HeapInuse)
	m.HeapObjects = Gauge(memStats.HeapObjects)
	m.HeapReleased = Gauge(memStats.HeapReleased)
	m.HeapSys = Gauge(memStats.HeapSys)
	m.LastGC = Gauge(memStats.LastGC)
	m.Lookups = Gauge(memStats.Lookups)
	m.MCacheInuse = Gauge(memStats.MCacheInuse)
	m.MCacheSys = Gauge(memStats.MCacheSys)
	m.MSpanInuse = Gauge(memStats.MSpanInuse)
	m.MSpanSys = Gauge(memStats.MSpanSys)
	m.Mallocs = Gauge(memStats.Mallocs)
	m.NextGC = Gauge(memStats.NextGC)
	m.NumForcedGC = Gauge(memStats.NumForcedGC)
	m.NumGC = Gauge(memStats.NumGC)
	m.OtherSys = Gauge(memStats.OtherSys)
	m.PauseTotalNs = Gauge(memStats.PauseTotalNs)
	m.StackInuse = Gauge(memStats.StackInuse)
	m.StackSys = Gauge(memStats.StackSys)
	m.Sys = Gauge(memStats.Sys)
}

func makeMetrics(pollCount Counter) Metrics {
	metrics := Metrics{
		PollCount:   pollCount,
		RandomValue: Gauge(rand.Float64()),
	}
	metrics.feelMemStats()
	return metrics
}

type MetricsConsumer interface {
	Consume(*Metrics)
}

type MetricsConsumerFunc func(*Metrics)

func (f MetricsConsumerFunc) Consume(metrics *Metrics) {
	log.Print("Metrics: Consume")
	f(metrics)
}
