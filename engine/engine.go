package engine

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/xh3b4sd/breakr"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/locker"
	"github.com/xh3b4sd/redigo/pool"
	"github.com/xh3b4sd/rescue/balancer"
	"github.com/xh3b4sd/rescue/metric"
	"github.com/xh3b4sd/rescue/timer"
	"github.com/xh3b4sd/tracer"
)

const (
	// Expiry is the default expiry of any given task.
	Expiry = 30 * time.Second
)

const (
	// Week is the time.Duration of 7 days.
	Week = 7 * 24 * time.Hour
)

type Config struct {
	Balancer balancer.Interface
	Cleanup  time.Duration
	Expiry   time.Duration
	Locker   locker.Interface
	Logger   logger.Interface
	Metric   *metric.Collection
	Queue    string
	Redigo   redigo.Interface
	Sepkey   string
	Timer    *timer.Timer
	Worker   string
}

type Engine struct {
	bal balancer.Interface
	// cac is the local lookup table for tasks that have been chosen to be
	// processed without assigning direct ownership to this particular worker
	// process. An example of necessary mappings we need to track for workers are
	// all tasks defining the delivery method "all".
	cac map[string]*local
	cln time.Duration
	ctx context.Context
	exp time.Duration
	loc locker.Interface
	log logger.Interface
	met *metric.Collection
	// pnt is the local point in time at which this worker became operational.
	// Further, this pointer will move forward with every broadcasted task that
	// got completed locally. This pointer will be used to e.g. decide whether to
	// process broadcasted tasks declared with method "all".
	pnt time.Time
	que string
	red redigo.Interface
	sep string
	tim *timer.Timer
	// wrk is the identifier of this worker process.
	wrk string
}

func New(config Config) *Engine {
	if config.Balancer == nil {
		config.Balancer = balancer.Default()
	}
	if config.Cleanup == 0 {
		config.Cleanup = Week
	}
	if config.Expiry == 0 {
		config.Expiry = Expiry
	}
	if config.Logger == nil {
		config.Logger = logger.Default()
	}
	if config.Metric == nil {
		config.Metric = metric.Default()
	}
	if config.Queue == "" {
		config.Queue = "default"
	}
	if config.Redigo == nil {
		config.Redigo = redigo.Default()
	}
	if config.Locker == nil {
		config.Locker = defLoc(config.Redigo.Listen())
	}
	if config.Sepkey == "" {
		config.Sepkey = ":"
	}
	if config.Timer == nil {
		config.Timer = timer.New()
	}
	if config.Worker == "" {
		config.Worker = uuid.New().String()
	}

	e := &Engine{
		bal: config.Balancer,
		cac: map[string]*local{},
		cln: config.Cleanup,
		ctx: context.Background(),
		exp: config.Expiry,
		loc: config.Locker,
		log: config.Logger,
		met: config.Metric,
		pnt: config.Timer.Engine(),
		que: config.Queue,
		red: config.Redigo,
		sep: config.Sepkey,
		tim: config.Timer,
		wrk: config.Worker,
	}

	return e
}

func (e *Engine) lerror(err error) {
	e.log.Log(
		e.ctx,
		"level", "error",
		"message", err.Error(),
		"stack", tracer.Stack(err),
	)
}

func defLoc(add string) locker.Interface {
	return locker.New(locker.Config{
		Brk: breakr.New(breakr.Config{
			Failure: breakr.Failure{
				Budget: 30,
				Cooler: 1 * time.Second,
			},
		}),
		Poo: pool.NewSinglePoolWithAddress(add),
	})
}
