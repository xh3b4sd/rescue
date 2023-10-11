package rescue

import "github.com/xh3b4sd/rescue/engine"

func Default() Interface {
	return engine.New(engine.Config{})
}

func Fake() Interface {
	return engine.New(engine.Config{})
}
