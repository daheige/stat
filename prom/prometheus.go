// Package prom for golang prometheus
package prom

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// LibClient for mc,redis,db(mysql)
	LibClient = New().WithTimer("go_lib_client", []string{"method", "name"}).
			WithState("go_lib_client_state", []string{"method", "name"}).
			WithCounter("go_lib_client_code", []string{"method", "code"})

	// RPCClient rpc client
	RPCClient = New().WithTimer("go_rpc_client", []string{"method", "name"}).
			WithState("go_rpc_client_state", []string{"method", "name"}).
			WithCounter("go_rpc_client_code", []string{"method", "name", "code"})

	// RPCServer for rpc server
	RPCServer = New().WithTimer("go_rpc_server", []string{"method", "name"}).
			WithCounter("go_rpc_server_code", []string{"method", "name", "code"})

	// HTTPClient http client
	HTTPClient = New().WithTimer("go_http_client", []string{"method", "name"}).
			WithState("go_http_client_state", []string{"method", "name"}).
			WithCounter("go_http_client_code", []string{"method", "name", "code"})

	// HTTPServer for http server
	HTTPServer = New().WithTimer("go_http_server", []string{"method", "name"}).
			WithCounter("go_http_server_code", []string{"method", "name", "code"})

	// APPErrorCount for business err count
	// RBI monitoring for business errors
	APPErrorCount = New().WithCounter("go_app_error_count", []string{"name"}).
			WithState("go_app_error_state", []string{"name"})

	// APPInfoCount for business info count
	// routine monitoring of normal business
	APPInfoCount = New().WithCounter("go_app_info_count", []string{"name"}).
			WithState("go_app_info_state", []string{"name"})

	// CacheHit for cache hit
	CacheHit = New().WithCounter("go_cache_hit", []string{"name"})

	// CacheMiss for cache miss
	CacheMiss = New().WithCounter("go_cache_miss", []string{"name"})

	// DBQuery db query time and state
	DBQuery = New().WithTimer("go_db_query", []string{"method", "name"}).
		WithSummary("go_db_query_summary", []string{"method", "name"}).
		WithState("go_db_query_state", []string{"method", "name"})
)

// Prom struct.
type Prom struct {
	timer        *prometheus.HistogramVec
	timerSummary *prometheus.SummaryVec
	counter      *prometheus.CounterVec
	state        *prometheus.GaugeVec
}

// New creates a Prom instance.
func New() *Prom {
	return &Prom{}
}

// WithTimer with summary timer
func (p *Prom) WithTimer(name string, labels []string) *Prom {
	if p == nil || p.timer != nil {
		return p
	}

	p.timer = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name,
			Help: name,
		}, labels)

	prometheus.MustRegister(p.timer)
	return p
}

// WithSummary with summary timer
func (p *Prom) WithSummary(name string, labels []string) *Prom {
	if p == nil || p.timerSummary != nil {
		return p
	}

	p.timerSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: name,
			Help: name,
		}, labels)

	prometheus.MustRegister(p.timerSummary)

	return p
}

// WithCounter sets counter.
func (p *Prom) WithCounter(name string, labels []string) *Prom {
	if p == nil || p.counter != nil {
		return p
	}

	p.counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: name,
		}, labels)

	prometheus.MustRegister(p.counter)

	return p
}

// WithState sets state.
func (p *Prom) WithState(name string, labels []string) *Prom {
	if p == nil || p.state != nil {
		return p
	}

	p.state = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: name,
		}, labels)

	prometheus.MustRegister(p.state)

	return p
}

// Timing log timing information (in milliseconds) without sampling
func (p *Prom) Timing(name string, time int64, extra ...string) {
	label := append([]string{name}, extra...)
	if p.timer != nil {
		p.timer.WithLabelValues(label...).Observe(float64(time))
	}
}

// Incr increments one stat counter without sampling
func (p *Prom) Incr(name string, extra ...string) {
	label := append([]string{name}, extra...)
	if p.counter != nil {
		p.counter.WithLabelValues(label...).Inc()
	}

	if p.state != nil {
		p.state.WithLabelValues(label...).Inc()
	}
}

// Decr decrements one stat counter without sampling
func (p *Prom) Decr(name string, extra ...string) {
	if p.state != nil {
		label := append([]string{name}, extra...)
		p.state.WithLabelValues(label...).Dec()
	}
}

// State set state
func (p *Prom) State(name string, v int64, extra ...string) {
	if p.state != nil {
		label := append([]string{name}, extra...)
		p.state.WithLabelValues(label...).Set(float64(v))
	}
}

// Add add count v must > 0
func (p *Prom) Add(name string, v int64, extra ...string) {
	label := append([]string{name}, extra...)
	if p.counter != nil {
		p.counter.WithLabelValues(label...).Add(float64(v))
	}

	if p.state != nil {
		p.state.WithLabelValues(label...).Add(float64(v))
	}
}
