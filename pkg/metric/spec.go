package metric

import "github.com/prometheus/client_golang/prometheus"

type Interface interface {
	Dec()
	Des() *prometheus.Desc
	Get() float64
	Inc()
	Res()
	Set(i float64)
	Sin(o func() error) error
}
