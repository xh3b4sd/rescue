package engine

import (
	"testing"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo/pkg/fake"
)

func Test_Engine_Interface(t *testing.T) {
	var e *Engine
	{
		e = New(Config{
			Logger: logger.Fake(),
			Redigo: fake.New(),
		})
	}

	var _ Interface = e
}
