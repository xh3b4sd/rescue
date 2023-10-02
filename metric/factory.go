package metric

import "github.com/prometheus/client_golang/prometheus"

func Default() *Collection {
	c := &Collection{
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
			Exists: &CollectionEngineCollector{
				Cal: &Metric{d: prometheus.NewDesc("rescue_engine_exists_call_total" /*********/, "the number of times a call to Engine.Exists was made", nil, nil)},
				Dur: &Metric{d: prometheus.NewDesc("rescue_engine_exists_duration_seconds" /***/, "the number of seconds a call to Engine.Exists took", nil, nil)},
				Err: &Metric{d: prometheus.NewDesc("rescue_engine_exists_error_total" /********/, "the number of errors a call to Engine.Exists produced", nil, nil)},
			},
			Expire: &CollectionEngineCollector{
				Cal: &Metric{d: prometheus.NewDesc("rescue_engine_expire_call_total" /*********/, "the number of times a call to Engine.Expire was made", nil, nil)},
				Dur: &Metric{d: prometheus.NewDesc("rescue_engine_expire_duration_seconds" /***/, "the number of seconds a call to Engine.Expire took", nil, nil)},
				Err: &Metric{d: prometheus.NewDesc("rescue_engine_expire_error_total" /********/, "the number of errors a call to Engine.Expire produced", nil, nil)},
			},
			Extend: &CollectionEngineCollector{
				Cal: &Metric{d: prometheus.NewDesc("rescue_engine_extend_call_total" /*********/, "the number of times a call to Engine.Extend was made", nil, nil)},
				Dur: &Metric{d: prometheus.NewDesc("rescue_engine_extend_duration_seconds" /***/, "the number of seconds a call to Engine.Extend took", nil, nil)},
				Err: &Metric{d: prometheus.NewDesc("rescue_engine_extend_error_total" /********/, "the number of errors a call to Engine.Extend produced", nil, nil)},
			},
			Lister: &CollectionEngineCollector{
				Cal: &Metric{d: prometheus.NewDesc("rescue_engine_lister_call_total" /*********/, "the number of times a call to Engine.Lister was made", nil, nil)},
				Dur: &Metric{d: prometheus.NewDesc("rescue_engine_lister_duration_seconds" /***/, "the number of seconds a call to Engine.Lister took", nil, nil)},
				Err: &Metric{d: prometheus.NewDesc("rescue_engine_lister_error_total" /********/, "the number of errors a call to Engine.Lister produced", nil, nil)},
			},
			Search: &CollectionEngineCollector{
				Cal: &Metric{d: prometheus.NewDesc("rescue_engine_search_call_total" /*********/, "the number of times a call to Engine.Search was made", nil, nil)},
				Dur: &Metric{d: prometheus.NewDesc("rescue_engine_search_duration_seconds" /***/, "the number of seconds a call to Engine.Search took", nil, nil)},
				Err: &Metric{d: prometheus.NewDesc("rescue_engine_search_error_total" /********/, "the number of errors a call to Engine.Search produced", nil, nil)},
			},
		},
		Task: &CollectionTask{
			Expired:  &Metric{d: prometheus.NewDesc("rescue_task_expired_total" /****/, "the number of times a task was expired during a call to Engine.Expire", nil, nil)},
			Extended: &Metric{d: prometheus.NewDesc("rescue_task_extended_total" /***/, "the number of times a task was extended during a call to Engine.Extend", nil, nil)},
			Inactive: &Metric{d: prometheus.NewDesc("rescue_task_inactive_total" /***/, "the number of tasks found idle in the queue during a call to Engine.*", nil, nil)},
			NotFound: &Metric{d: prometheus.NewDesc("rescue_task_notfound_total" /***/, "the number of times a task could not be found during a call to Engine.*", nil, nil)},
			Obsolete: &Metric{d: prometheus.NewDesc("rescue_task_obsolete_total" /***/, "the number of times a nested task was removed during a call to Engine.Search", nil, nil)},
			Outdated: &Metric{d: prometheus.NewDesc("rescue_task_outdated_total" /***/, "the number of times a task has changed internally during a call to Engine.*", nil, nil)},
			Parallel: &Metric{d: prometheus.NewDesc("rescue_task_parallel_total" /***/, "the number of tasks claimed concurrently by a single worker during a call to Engine.Search", nil, nil)},
		},
	}

	c.Reset()

	return c
}
