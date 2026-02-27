package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Business metrics
	IncidentsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "incidents_total",
			Help: "Total number of incidents created",
		},
		[]string{"severity"},
	)

	IncidentsByStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "incidents_by_status",
			Help: "Current number of incidents by status",
		},
		[]string{"status"},
	)

	IncidentResolutionTime = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "incident_resolution_time_seconds",
			Help:    "Time taken to resolve incidents in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	IncidentsAssignedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "incidents_assigned_total",
			Help: "Total number of incident assignments",
		},
	)

	// Technical metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"method", "endpoint"},
	)
)

// RecordIncidentCreated registra la creación de un incidente
func RecordIncidentCreated(severity string) {
	IncidentsTotal.WithLabelValues(severity).Inc()
}

// RecordIncidentAssigned registra la asignación de un incidente
func RecordIncidentAssigned() {
	IncidentsAssignedTotal.Inc()
}

// RecordIncidentResolved registra la resolución de un incidente
func RecordIncidentResolved(duration time.Duration) {
	IncidentResolutionTime.Observe(duration.Seconds())
}

// RecordHTTPRequest registra una petición HTTP
func RecordHTTPRequest(method, endpoint, status string, duration time.Duration) {
	HTTPRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}
