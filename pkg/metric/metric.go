package metric

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xh3b4sd/tracer"
)

type Metric struct {
	d *prometheus.Desc
	i float64
	m sync.Mutex
}

func (m *Metric) Dec() {
	m.m.Lock()
	defer m.m.Unlock()

	m.i = m.i - 1
}

func (m *Metric) Des() *prometheus.Desc {
	return m.d
}

func (m *Metric) Get() float64 {
	m.m.Lock()
	defer m.m.Unlock()

	return m.i
}

func (m *Metric) Inc() {
	m.m.Lock()
	defer m.m.Unlock()

	m.i = m.i + 1
}

func (m *Metric) Set(i float64) {
	m.m.Lock()
	defer m.m.Unlock()

	m.i = i
}

func (m *Metric) Sin(o func() error) error {
	m.m.Lock()
	defer m.m.Unlock()

	s := time.Now()
	defer func() {
		m.i = time.Since(s).Seconds()
	}()

	err := o()
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
