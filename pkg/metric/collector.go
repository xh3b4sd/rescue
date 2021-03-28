package metric

import (
	"sync"
	"time"

	"github.com/xh3b4sd/tracer"
)

type collector struct {
	i float64
	m sync.Mutex
}

func (c *collector) Dec() {
	c.m.Lock()
	defer c.m.Unlock()

	c.i = c.i - 1
}

func (c *collector) Dur(o func() error) error {
	c.m.Lock()
	defer c.m.Unlock()

	s := time.Now()
	defer func() {
		c.i = time.Since(s).Seconds()
	}()

	err := o()
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (c *collector) Get() float64 {
	c.m.Lock()
	defer c.m.Unlock()

	return c.i
}

func (c *collector) Inc() {
	c.m.Lock()
	defer c.m.Unlock()

	c.i = c.i + 1
}

func (c *collector) Set(i float64) {
	c.m.Lock()
	defer c.m.Unlock()

	c.i = i
}
