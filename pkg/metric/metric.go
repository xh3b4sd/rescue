package metric

type Metric struct {
	Engine *MetricEngine
	Task   *MetricTask
}

type MetricEngine struct {
	Create *MetricEngineCollector
	Delete *MetricEngineCollector
	Expire *MetricEngineCollector
	Metric *MetricEngineCollector
	Search *MetricEngineCollector
}

type MetricEngineCollector struct {
	Action Collector
	Error  Collector
}

type MetricTask struct {
	Expired  Collector
	NotFound Collector
	Outdated Collector
}

func New() *Metric {
	m := &Metric{
		Engine: &MetricEngine{
			Create: &MetricEngineCollector{
				Action: &collector{},
				Error:  &collector{},
			},
			Delete: &MetricEngineCollector{
				Action: &collector{},
				Error:  &collector{},
			},
			Expire: &MetricEngineCollector{
				Action: &collector{},
				Error:  &collector{},
			},
			Metric: &MetricEngineCollector{
				Action: &collector{},
				Error:  &collector{},
			},
			Search: &MetricEngineCollector{
				Action: &collector{},
				Error:  &collector{},
			},
		},
		Task: &MetricTask{
			Expired:  &collector{},
			NotFound: &collector{},
			Outdated: &collector{},
		},
	}

	return m
}

func (m *Metric) Reset() {
	m.Engine.Create.Action.Set(0)
	m.Engine.Create.Error.Set(0)

	m.Engine.Delete.Action.Set(0)
	m.Engine.Delete.Error.Set(0)

	m.Engine.Expire.Action.Set(0)
	m.Engine.Expire.Error.Set(0)

	m.Engine.Metric.Action.Set(0)
	m.Engine.Metric.Error.Set(0)

	m.Engine.Search.Action.Set(0)
	m.Engine.Search.Error.Set(0)

	m.Task.Expired.Set(0)
	m.Task.NotFound.Set(0)
	m.Task.Outdated.Set(0)
}
