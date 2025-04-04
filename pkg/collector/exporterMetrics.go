package collector

import "github.com/prometheus/client_golang/prometheus"

// ExporterMetrics tracks metrics internal to the exporter
type ExporterMetrics struct {
	refreshErrors *prometheus.CounterVec
}

func NewExporterMetrics() *ExporterMetrics {
	return &ExporterMetrics{
		refreshErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "exporter",
			Name:      "data_model_refresh_errors_total",
			Help:      "Errors encountered while updating the internal server model",
		}, []string{"collector"}),
	}
}

// Init initializes the metrics given a collector label to make sure that they appear even though they might be never updated
func (i *ExporterMetrics) Init(collector string) {
	i.refreshErrors.WithLabelValues(collector)
}

// RefreshError increases the refresh error counter of the given collector by one.
func (i *ExporterMetrics) RefreshError(collector string) {
	i.refreshErrors.WithLabelValues(collector).Inc()
}

func (i *ExporterMetrics) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(i, descs)
}

func (i *ExporterMetrics) Collect(metrics chan<- prometheus.Metric) {
	// Ensure that the metric is always initialized
	i.refreshErrors.Collect(metrics)
}
