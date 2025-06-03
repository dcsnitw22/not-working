// Add these metrics in amf/metrics/metrics.go

package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Mutex for metrics synchronization
	metricsMutex sync.RWMutex

	// Registration metrics
	RegistrationAttempts = promauto.NewCounter(prometheus.CounterOpts{
		Name: "amf_registration_operations_total",
		Help: "Total number of registration attempts",
	})

	RegistrationSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "amf_registration_operations_successful",
		Help: "Total successful registrations",
	})

	RegistrationFailures = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "amf_registration_operations_failed",
		Help: "Registration failures by cause",
	}, []string{"cause"})

	CreateAttempts = promauto.NewCounter(prometheus.CounterOpts{
		Name: "amf_session_create_operations_total",
		Help: "Total number of registration attempts",
	})

	CreateSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "amf_session_create_operations_successful",
		Help: "Total successful registrations",
	})

	CreateFailures = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "amf_session_create_operations_failure",
		Help: "Registration failures by cause",
	}, []string{"cause"})

	// UE Context metrics
	UeContextGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "amf_active_ue_contexts",
		Help: "Current number of active UE contexts",
	})

	// Procedure timing
	RegistrationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "amf_registration_duration_seconds",
		Help:    "Time taken for registration procedures",
		Buckets: prometheus.ExponentialBuckets(0.1, 2, 5),
	}, []string{"outcome"})

	// NGAP Connection metrics
	NgapConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "amf_ngap_connections",
		Help: "Current active NGAP connections",
	})

	// AMF State metrics
	AmfStateGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "amf_state",
		Help: "Current state of AMF procedures",
	}, []string{"state_type"})

	AmfcreateProcess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "amf_create_operations_total",
		Help: "The total number of AMF create operations",
	})

	AmfupdateProcess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "amf_update_operations_total",
		Help: "The total number of AMF update operations",
	})

	AmfretrieveProcess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "amf_retrieve_operations_total",
		Help: "The total number of AMF retrieve operations",
	})

	AmfreleaseProcess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "amf_release_operations_total",
		Help: "The total number of AMF release operations",
	})

	AmfcreateSessionSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "amf_create_session_success_total",
		Help: "The total number of successful session creations",
	})

	AmfcreateSessionAttempts = promauto.NewCounter(prometheus.CounterOpts{
		Name: "amf_create_session_attempts_total",
		Help: "The total number of session creation attempts",
	})
)

// Thread-safe metric update functions
func UpdateCounter(counter prometheus.Counter) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()
	counter.Inc()
}

func UpdateCounterVec(counterVec *prometheus.CounterVec, value float64, labels ...string) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()
	counterVec.WithLabelValues(labels...).Add(value)
}

func UpdateGauge(gauge prometheus.Gauge, value float64) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()
	gauge.Set(value)
}

func UpdateGaugeVec(gaugeVec *prometheus.GaugeVec, value float64, labels ...string) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()
	gaugeVec.WithLabelValues(labels...).Set(value)
}

func UpdateHistogram(histogram *prometheus.HistogramVec, value float64, labels ...string) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()
	histogram.WithLabelValues(labels...).Observe(value)
}
