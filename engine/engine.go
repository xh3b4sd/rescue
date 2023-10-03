package engine

import (
	"context"
	"time"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/rescue/balancer"
	"github.com/xh3b4sd/rescue/metric"
	"github.com/xh3b4sd/rescue/random"
	"github.com/xh3b4sd/tracer"
)

const (
	// Expiry is the default expiry of any given task.
	Expiry = 30 * time.Second
)

type Config struct {
	Balancer balancer.Interface
	Expiry   time.Duration
	Logger   logger.Interface
	Metric   *metric.Collection
	Queue    string
	Redigo   redigo.Interface
	Worker   string
}

type Engine struct {
	bal balancer.Interface
	ctx context.Context
	exp time.Duration
	log logger.Interface
	met *metric.Collection
	que string
	red redigo.Interface
	wrk string
}

func New(config Config) *Engine {
	if config.Balancer == nil {
		config.Balancer = balancer.Default()
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
	if config.Worker == "" {
		config.Worker = random.New()
	}

	e := &Engine{
		bal: config.Balancer,
		ctx: context.Background(),
		exp: config.Expiry,
		log: config.Logger,
		met: config.Metric,
		que: config.Queue,
		red: config.Redigo,
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
