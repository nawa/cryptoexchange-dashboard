package prometheus

import (
	"net/http"
	"time"

	"github.com/kataras/iris/context"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// DefaultBuckets prometheus buckets.
	DefaultBuckets = []float64{300, 1200, 5000}
)

const (
	reqsName    = "requests_total"
	latencyName = "request_duration_milliseconds"
)

// Prometheus is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
//
// Usage: pass its `ServeHTTP` to a route or globally.
type Prometheus struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

// New returns a new prometheus middleware.
//
// If buckets are empty then `DefaultBuckets` are setted.
func New(name string, buckets ...float64) *Prometheus {
	p := Prometheus{}
	p.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(p.reqs)

	if len(buckets) == 0 {
		buckets = DefaultBuckets
	}

	p.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(p.latency)

	return &p
}

func (p *Prometheus) ServeHTTP(ctx context.Context) {
	start := time.Now()
	ctx.Next()
	r := ctx.Request()

	p.reqs.WithLabelValues(http.StatusText(ctx.GetStatusCode()), r.Method, r.URL.Path).
		Inc()

	p.latency.WithLabelValues(http.StatusText(ctx.GetStatusCode()), r.Method, r.URL.Path).
		Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
}
