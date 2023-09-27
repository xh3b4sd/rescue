package engine

import (
	"testing"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo/pkg/fake"

	"github.com/xh3b4sd/rescue"
)

func Test_Engine_Interface(t *testing.T) {
	var err error

	var e *Engine
	{
		c := Config{
			Logger: logger.Fake(),
			Redigo: fake.New(),
		}

		e, err = New(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	var _ rescue.Interface = e
}
