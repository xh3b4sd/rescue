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
	Metric *metric.Collection
	Redigo redigo.Interface

	Owner string
	TTL   time.Duration
}

type Engine struct {
	logger logger.Interface
	metric *metric.Collection
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

	e.metric.Engine.Create.Cal.Inc()

	o := func() error {
		err = e.create(tsk)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.metric.Engine.Create.Dur.Sin(o)
	if err != nil {
		e.metric.Engine.Create.Err.Inc()
		return tracer.Mask(err)
	}

	return nil
}

func (e *Engine) Delete(tsk *task.Task) error {
	var err error

	e.metric.Engine.Delete.Cal.Inc()

	o := func() error {
		err = e.delete(tsk)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.metric.Engine.Delete.Dur.Sin(o)
	if err != nil {
		e.metric.Engine.Delete.Err.Inc()
		return tracer.Mask(err)
	}

	return nil
}

func (e *Engine) Exists(tsk *task.Task) (bool, error) {
	var err error
	var exi bool

	e.metric.Engine.Exists.Cal.Inc()

	o := func() error {
		exi, err = e.exists(tsk)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.metric.Engine.Exists.Dur.Sin(o)
	if err != nil {
		e.metric.Engine.Delete.Err.Inc()
		return false, tracer.Mask(err)
	}

	return exi, nil
}

func (e *Engine) Expire() error {
	var err error

	e.metric.Engine.Expire.Cal.Inc()

	o := func() error {
		err = e.expire()
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.metric.Engine.Expire.Dur.Sin(o)
	if err != nil {
		e.metric.Engine.Expire.Err.Inc()
		return tracer.Mask(err)
	}

	return nil
}

func (e *Engine) Search() (*task.Task, error) {
	var err error
	var tsk *task.Task

	e.metric.Engine.Search.Cal.Inc()

	o := func() error {
		tsk, err = e.search()
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.metric.Engine.Search.Dur.Sin(o)
	if err != nil {
		e.metric.Engine.Search.Err.Inc()
		return nil, tracer.Mask(err)
	}

	return tsk, nil
}
