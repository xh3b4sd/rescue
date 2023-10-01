package engine

import (
	"testing"
)

func Test_Engine_Interface(t *testing.T) {
	var _ Interface = New(Config{})
}
