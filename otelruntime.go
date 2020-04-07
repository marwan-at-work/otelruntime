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

// tag keys
const (
	Goroutines = "cpu_goroutines"
	Heap       = "heap_alloc"
	Pause      = "pause_ns"
	SysHeap    = "sys_heap"
)

// Config for the runtime reporter
type Config struct {
	meter metric.Meter
	name  string
}

// Option allows changes to the Config
type Option func(c *Config)

// WithMeter sets the meter that the runtime reporter will use
// otherwise global.Meter is used by default.
func WithMeter(meter metric.Meter) Option {
	return func(c *Config) {
		c.meter = meter
	}
}

// WithMetricName customizes the exported metric name or otherwise
// will be called go_runtime by default.
func WithMetricName(name string) Option {
	return func(c *Config) {
		c.name = name
	}
}

// Register registers an infinite loop that
// profiles the process and reports the views
func Register(options ...Option) {
	c := &Config{meter: global.Meter("otelruntime"), name: "go_runtime"}
	for _, o := range options {
		o(c)
	}
	mm := metric.Must(c.meter)
	mm.RegisterInt64Observer(c.name, func(result metric.Int64ObserverResult) {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		result.Observe(int64(runtime.NumCPU()), ProcessKey.String(Goroutines))
		result.Observe(int64(ms.HeapAlloc), ProcessKey.String(Heap))
		result.Observe(int64(ms.PauseNs[(ms.NumGC+255)%256]), ProcessKey.String(Pause))
		result.Observe(int64(ms.Sys), ProcessKey.String(SysHeap))
	})
}
