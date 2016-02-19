package metrics

import (
	"sync"

	goMetrics "github.com/rcrowley/go-metrics"
)

const (
	KgMetrics = iota
)

type MetricsFactory struct {
}

func (mFactory MetricsFactory) GetMetrics(metricsType int) Metrics {
	switch metricsType {
	case KgMetrics:
		return NewKgMetrics()
	default:
		return nil
	}
}

type Metrics interface {
	Mark(string, int64)
}

// Metrics ...
type mMetrics struct {
	registry *goMetrics.Registry
}

var singleM Metrics
var onceLock sync.Once

// NewMetrics ...
func NewKgMetrics() Metrics {
	onceLock.Do(
		func() {
			if singleM == nil {
				singleM = &mMetrics{
					registry: &goMetrics.DefaultRegistry,
				}
			}
		})
	return singleM
}

// Mark ...
func (m *mMetrics) Mark(meterName string, n int64) {
	meter := goMetrics.GetOrRegisterMeter(meterName, *m.registry)
	meter.Mark(n)
}
