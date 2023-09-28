package engine

import (
	"testing"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo/pkg/fake"

	"github.com/xh3b4sd/rescue"
)

func Test_Engine_Interface(t *testing.T) {
	var e *Engine
	{
		e = New(Config{
			Logger: logger.Fake(),
			Redigo: fake.New(),
		})
	}

	var _ rescue.Interface = e
}
