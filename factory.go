package rescue

import (
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/locker"
	"github.com/xh3b4sd/rescue/engine"
)

func Default() Interface {
	return engine.New(engine.Config{
		Redigo: redigo.Default(),
	})
}

func Fake() Interface {
	return engine.New(engine.Config{
		Locker: &locker.Fake{},
		Redigo: redigo.Fake(),
	})
}
