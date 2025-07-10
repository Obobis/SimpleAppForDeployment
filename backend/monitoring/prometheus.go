package monitoring

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics содержит все метрики приложения
type Metrics struct {
	WSConnections   prometheus.Gauge
	DBQueryDuration prometheus.Histogram
	HTTPRequests    *prometheus.CounterVec
}

// NewMetrics создает и регистрирует метрики
func NewMetrics() *Metrics {
	m := &Metrics{
		WSConnections: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "websocket_active_connections",
			Help: "Current number of active WebSocket connections",
		}),

		DBQueryDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration distribution",
			Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1},
		}),

		HTTPRequests: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		}, []string{"method", "path", "status"}),
	}

	return m
}

// Handler возвращает HTTP handler для метрик
func (m *Metrics) Handler() http.Handler {
	return promhttp.Handler()
}

// InstrumentHTTP добавляет мониторинг HTTP запросов
func (m *Metrics) InstrumentHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логика для отслеживания запросов
		recorder := &StatusRecorder{w, 200}
		next.ServeHTTP(recorder, r)

		m.HTTPRequests.WithLabelValues(
			r.Method,
			r.URL.Path,
			fmt.Sprintf("%d", recorder.Status),
		).Inc()
	})
}

// StatusRecorder перехватывает статус ответа
type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}
