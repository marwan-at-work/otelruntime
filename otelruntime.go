package otelruntime

import (
	"runtime"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
)

// ProcessKey is the key that differentiates
// the metric type because only one metric is exported
// for all stats. Use this key to switch between the metrics
// below
var ProcessKey = key.New("go_process")

const (
	Goroutines = "cpu_goroutines"
	Heap       = "heap_alloc"
	Pause      = "pause_ns"
	SysHeap    = "sys_heap"
)

// Register registers an infinite loop that
// profiles the process and reports the views
func Register() {
	meter := global.Meter("otelruntime")
	mm := metric.Must(meter)
	mm.RegisterInt64Observer("go_process", func(result metric.Int64ObserverResult) {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		result.Observe(int64(runtime.NumCPU()), ProcessKey.String(Goroutines))
		result.Observe(int64(ms.HeapAlloc), ProcessKey.String(Heap))
		result.Observe(int64(ms.PauseNs[(ms.NumGC+255)%256]), ProcessKey.String(Pause))
		result.Observe(int64(ms.Sys), ProcessKey.String(SysHeap))
	})
}
