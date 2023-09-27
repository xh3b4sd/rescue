package engine

import (
	"context"
	"time"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/balancer"
	"github.com/xh3b4sd/rescue/pkg/metric"
	"github.com/xh3b4sd/rescue/pkg/random"
)

const (
	// TTL is the default expiry of any given task.
	TTL = 30 * time.Second
)

type Config struct {
	Balancer balancer.Interface
	Logger   logger.Interface
	Metric   *metric.Collection
	Owner    string
	Queue    string
	Redigo   redigo.Interface
	TTL      time.Duration
}

type Engine struct {
	bal balancer.Interface
	ctx context.Context
	log logger.Interface
	met *metric.Collection
	own string
	que string
	red redigo.Interface
	ttl time.Duration
}

func New(config Config) (*Engine, error) {
	if config.Balancer == nil {
		config.Balancer = balancer.Default()
	}
	if config.Logger == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Metric == nil {
		config.Metric = metric.Default()
	}
	if config.Owner == "" {
		config.Owner = random.MustNew()
	}
	if config.Queue == "" {
		config.Queue = "def"
	}
	if config.Redigo == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Redigo must not be empty", config)
	}
	if config.TTL == 0 {
		config.TTL = TTL
	}

	e := &Engine{
		bal: config.Balancer,
		ctx: context.Background(),
		log: config.Logger,
		met: config.Metric,
		own: config.Owner,
		que: config.Queue,
		red: config.Redigo,
		ttl: config.TTL,
	}

	return e, nil
}
