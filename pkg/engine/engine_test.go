package engine

import (
	"testing"

	loggerfake "github.com/xh3b4sd/logger/fake"
	redigofake "github.com/xh3b4sd/redigo/fake"

	"github.com/xh3b4sd/rescue"
)

func Test_Engine_Interface(t *testing.T) {
	var err error

	var e *Engine
	{
		c := Config{
			Logger: loggerfake.New(),
			Redigo: redigofake.New(),
		}

		e, err = New(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	var _ rescue.Interface = e
}
