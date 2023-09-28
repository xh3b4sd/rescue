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
	s []float64
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

func (m *Metric) Res() {
	m.i = 0
	m.s = []float64{}
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
		m.s = append(m.s, time.Since(s).Seconds())
		m.i = avg(m.s)
	}()

	err := o()
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func avg(l []float64) float64 {
	var s float64
	for _, f := range l {
		s += f
	}

	return s / float64(len(l))
}
