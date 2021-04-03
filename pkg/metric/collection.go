package metric

import "github.com/prometheus/client_golang/prometheus"

type Collection struct {
	Engine *CollectionEngine
	Task   *CollectionTask
}

type CollectionEngine struct {
	Create *CollectionEngineCollector
	Delete *CollectionEngineCollector
	Expire *CollectionEngineCollector
	Metric *CollectionEngineCollector
	Search *CollectionEngineCollector
}

type CollectionEngineCollector struct {
	Cal Interface
	Dur Interface
	Err Interface
}

type CollectionTask struct {
	Expired  Interface
	NotFound Interface
	Outdated Interface
	Queued   Interface
}

func New() *Collection {
	m := &Collection{
		Engine: &CollectionEngine{
			Create: &CollectionEngineCollector{
				Cal: &Metric{d: prometheus.NewDesc("rescue_engine_create_call_total" /*********/, "the number of times a call to Engine.Create was made", nil, nil)},
				Dur: &Metric{d: prometheus.NewDesc("rescue_engine_create_duration_seconds" /***/, "the number of seconds a call to Engine.Create took", nil, nil)},
				Err: &Metric{d: prometheus.NewDesc("rescue_engine_create_error_total" /********/, "the number of errors a call to Engine.Create produced", nil, nil)},
			},
			Delete: &CollectionEngineCollector{
				Cal: &Metric{d: prometheus.NewDesc("rescue_engine_delete_call_total" /*********/, "the number of times a call to Engine.Delete was made", nil, nil)},
				Dur: &Metric{d: prometheus.NewDesc("rescue_engine_delete_duration_seconds" /***/, "the number of seconds a call to Engine.Delete took", nil, nil)},
				Err: &Metric{d: prometheus.NewDesc("rescue_engine_delete_error_total" /********/, "the number of errors a call to Engine.Delete produced", nil, nil)},
			},
			Expire: &CollectionEngineCollector{
				Cal: &Metric{d: prometheus.NewDesc("rescue_engine_expire_call_total" /*********/, "the number of times a call to Engine.Expire was made", nil, nil)},
				Dur: &Metric{d: prometheus.NewDesc("rescue_engine_expire_duration_seconds" /***/, "the number of seconds a call to Engine.Expire took", nil, nil)},
				Err: &Metric{d: prometheus.NewDesc("rescue_engine_expire_error_total" /********/, "the number of errors a call to Engine.Expire produced", nil, nil)},
			},
			Metric: &CollectionEngineCollector{
				Cal: &Metric{d: prometheus.NewDesc("rescue_engine_metric_call_total" /*********/, "the number of times a call to Engine.Metric was made", nil, nil)},
				Dur: &Metric{d: prometheus.NewDesc("rescue_engine_metric_duration_seconds" /***/, "the number of seconds a call to Engine.Metric took", nil, nil)},
				Err: &Metric{d: prometheus.NewDesc("rescue_engine_metric_error_total" /********/, "the number of errors a call to Engine.Metric produced", nil, nil)},
			},
			Search: &CollectionEngineCollector{
				Cal: &Metric{d: prometheus.NewDesc("rescue_engine_search_call_total" /*********/, "the number of times a call to Engine.Search was made", nil, nil)},
				Dur: &Metric{d: prometheus.NewDesc("rescue_engine_search_duration_seconds" /***/, "the number of seconds a call to Engine.Search took", nil, nil)},
				Err: &Metric{d: prometheus.NewDesc("rescue_engine_search_error_total" /********/, "the number of errors a call to Engine.Search produced", nil, nil)},
			},
		},
		Task: &CollectionTask{
			Expired:  &Metric{d: prometheus.NewDesc("rescue_task_expired_total" /****/, "the number of times a task was expired during a call to Engine.Expire", nil, nil)},
			NotFound: &Metric{d: prometheus.NewDesc("rescue_task_notfound_total" /***/, "the number of times a task was could not be found during a call to Engine.Search", nil, nil)},
			Outdated: &Metric{d: prometheus.NewDesc("rescue_task_outdated_total" /***/, "the number of times a task was tried to cleaned up during a call to Engine.Delete", nil, nil)},
			Queued:   &Metric{d: prometheus.NewDesc("rescue_task_queued_total" /*****/, "the number of tasks found in the queue during a call to Engine.Search", nil, nil)},
		},
	}

	return m
}

func (m *Collection) Reset() {
	m.Engine.Create.Cal.Set(0)
	m.Engine.Create.Dur.Set(0)
	m.Engine.Create.Err.Set(0)

	m.Engine.Delete.Cal.Set(0)
	m.Engine.Delete.Dur.Set(0)
	m.Engine.Delete.Err.Set(0)

	m.Engine.Expire.Cal.Set(0)
	m.Engine.Expire.Dur.Set(0)
	m.Engine.Expire.Err.Set(0)

	m.Engine.Metric.Cal.Set(0)
	m.Engine.Metric.Dur.Set(0)
	m.Engine.Metric.Err.Set(0)

	m.Engine.Search.Cal.Set(0)
	m.Engine.Search.Dur.Set(0)
	m.Engine.Search.Err.Set(0)

	m.Task.Expired.Set(0)
	m.Task.NotFound.Set(0)
	m.Task.Outdated.Set(0)
	m.Task.Queued.Set(0)
}
