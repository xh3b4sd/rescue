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

	m.Reset()

	return m
}

func (m *Collection) Reset() {
	m.Engine.Create.Cal.Res()
	m.Engine.Create.Dur.Res()
	m.Engine.Create.Err.Res()

	m.Engine.Delete.Cal.Res()
	m.Engine.Delete.Dur.Res()
	m.Engine.Delete.Err.Res()

	m.Engine.Expire.Cal.Res()
	m.Engine.Expire.Dur.Res()
	m.Engine.Expire.Err.Res()

	m.Engine.Search.Cal.Res()
	m.Engine.Search.Dur.Res()
	m.Engine.Search.Err.Res()

	m.Task.Expired.Res()
	m.Task.NotFound.Res()
	m.Task.Outdated.Res()
	m.Task.Queued.Res()
}
