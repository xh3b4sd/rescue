package engine

import (
	"time"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/metric"
	"github.com/xh3b4sd/rescue/pkg/random"
	"github.com/xh3b4sd/rescue/pkg/task"
)

type Config struct {
	Logger logger.Interface
	Metric *metric.Metric
	Redigo redigo.Interface

	Owner string
	TTL   time.Duration
}

type Engine struct {
	logger logger.Interface
	metric *metric.Metric
	redigo redigo.Interface

	owner string
	ttl   time.Duration
}

func New(config Config) (*Engine, error) {
	if config.Logger == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Metric == nil {
		config.Metric = metric.New()
	}
	if config.Redigo == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Redigo must not be empty", config)
	}

	if config.Owner == "" {
		config.Owner = random.MustNew()
	}
	if config.TTL == 0 {
		config.TTL = 30 * time.Second
	}

	e := &Engine{
		logger: config.Logger,
		metric: config.Metric,
		redigo: config.Redigo,

		owner: config.Owner,
		ttl:   config.TTL,
	}

	return e, nil
}

func (e *Engine) Create(tsk *task.Task) error {
	var err error

	e.metric.Engine.Create.Action.Inc()

	o := func() error {
		err = e.create(tsk)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.metric.Engine.Create.Action.Dur(o)
	if err != nil {
		e.metric.Engine.Create.Error.Inc()
		return tracer.Mask(err)
	}

	return nil
}

func (e *Engine) Delete(tsk *task.Task) error {
	var err error

	e.metric.Engine.Delete.Action.Inc()

	o := func() error {
		err = e.delete(tsk)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.metric.Engine.Delete.Action.Dur(o)
	if err != nil {
		e.metric.Engine.Delete.Error.Inc()
		return tracer.Mask(err)
	}

	return nil
}

func (e *Engine) Expire() error {
	var err error

	e.metric.Engine.Expire.Action.Inc()

	o := func() error {
		err = e.expire()
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.metric.Engine.Expire.Action.Dur(o)
	if err != nil {
		e.metric.Engine.Expire.Error.Inc()
		return tracer.Mask(err)
	}

	return nil
}

func (e *Engine) Metric() metric.Metric {
	e.metric.Engine.Metric.Action.Inc()

	return metric.Metric{}
}

func (e *Engine) Search() (*task.Task, error) {
	var err error
	var tsk *task.Task

	e.metric.Engine.Search.Action.Inc()

	o := func() error {
		tsk, err = e.search()
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.metric.Engine.Search.Action.Dur(o)
	if err != nil {
		e.metric.Engine.Search.Error.Inc()
		return nil, tracer.Mask(err)
	}

	return tsk, nil
}
