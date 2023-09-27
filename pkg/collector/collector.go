package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/metric"
)

type Config struct {
	Logger logger.Interface
	Metric *metric.Collection
}

type Collector struct {
	logger logger.Interface
	metric *metric.Collection
}

func New(config Config) (*Collector, error) {
	if config.Logger == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Metric == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Metric must not be empty", config)
	}

	c := &Collector{
		logger: config.Logger,
		metric: config.Metric,
	}

	return c, nil
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Create.Cal.Des() /***/, prometheus.CounterValue /***/, c.metric.Engine.Create.Cal.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Create.Dur.Des() /***/, prometheus.GaugeValue /*****/, c.metric.Engine.Create.Dur.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Create.Err.Des() /***/, prometheus.CounterValue /***/, c.metric.Engine.Create.Err.Get())

	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Delete.Cal.Des() /***/, prometheus.CounterValue /***/, c.metric.Engine.Delete.Cal.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Delete.Dur.Des() /***/, prometheus.GaugeValue /*****/, c.metric.Engine.Delete.Dur.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Delete.Err.Des() /***/, prometheus.CounterValue /***/, c.metric.Engine.Delete.Err.Get())

	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Exists.Cal.Des() /***/, prometheus.CounterValue /***/, c.metric.Engine.Exists.Cal.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Exists.Dur.Des() /***/, prometheus.GaugeValue /*****/, c.metric.Engine.Exists.Dur.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Exists.Err.Des() /***/, prometheus.CounterValue /***/, c.metric.Engine.Exists.Err.Get())

	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Expire.Cal.Des() /***/, prometheus.CounterValue /***/, c.metric.Engine.Expire.Cal.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Expire.Dur.Des() /***/, prometheus.GaugeValue /*****/, c.metric.Engine.Expire.Dur.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Expire.Err.Des() /***/, prometheus.CounterValue /***/, c.metric.Engine.Expire.Err.Get())

	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Extend.Cal.Des() /***/, prometheus.CounterValue /***/, c.metric.Engine.Extend.Cal.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Extend.Dur.Des() /***/, prometheus.GaugeValue /*****/, c.metric.Engine.Extend.Dur.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Extend.Err.Des() /***/, prometheus.CounterValue /***/, c.metric.Engine.Extend.Err.Get())

	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Lister.Cal.Des() /***/, prometheus.CounterValue /***/, c.metric.Engine.Lister.Cal.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Lister.Dur.Des() /***/, prometheus.GaugeValue /*****/, c.metric.Engine.Lister.Dur.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Lister.Err.Des() /***/, prometheus.CounterValue /***/, c.metric.Engine.Lister.Err.Get())

	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Search.Cal.Des() /***/, prometheus.CounterValue /***/, c.metric.Engine.Search.Cal.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Search.Dur.Des() /***/, prometheus.GaugeValue /*****/, c.metric.Engine.Search.Dur.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Engine.Search.Err.Des() /***/, prometheus.CounterValue /***/, c.metric.Engine.Search.Err.Get())

	ch <- prometheus.MustNewConstMetric(c.metric.Task.Expired.Des() /********/, prometheus.CounterValue /***/, c.metric.Task.Expired.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Task.Extended.Des() /*******/, prometheus.CounterValue /***/, c.metric.Task.Extended.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Task.NotFound.Des() /*******/, prometheus.CounterValue /***/, c.metric.Task.NotFound.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Task.Outdated.Des() /*******/, prometheus.CounterValue /***/, c.metric.Task.Outdated.Get())
	ch <- prometheus.MustNewConstMetric(c.metric.Task.Queued.Des() /*********/, prometheus.CounterValue /***/, c.metric.Task.Queued.Get())

	c.metric.Reset()
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}
