package engine

import (
	"context"
	"time"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/rescue/balancer"
	"github.com/xh3b4sd/rescue/metric"
	"github.com/xh3b4sd/rescue/random"
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

func New(config Config) *Engine {
	if config.Balancer == nil {
		config.Balancer = balancer.Default()
	}
	if config.Logger == nil {
		config.Logger = logger.Default()
	}
	if config.Metric == nil {
		config.Metric = metric.Default()
	}
	if config.Owner == "" {
		config.Owner = random.New()
	}
	if config.Queue == "" {
		config.Queue = "def"
	}
	if config.Redigo == nil {
		config.Redigo = redigo.Default()
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

	return e
}
