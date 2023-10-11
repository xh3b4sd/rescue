package rescue

import (
	"testing"
)

func Test_Factory_Interface_Default(t *testing.T) {
	var _ Interface = Default()
}
