//go:build redis

package conformance

import (
	"time"

	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/tracer"
)

func musTim(str string) time.Time {
	tim, err := time.Parse("2006-01-02T15:04:05.999999Z", str)
	if err != nil {
		panic(err)
	}

	return tim
}

// prgAll is a convenience function for calling FLUSHALL. The provided redigo
// interface is returned as is.
func prgAll(red redigo.Interface) redigo.Interface {
	{
		err := red.Purge()
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	return red
}
