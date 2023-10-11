package rescue

import (
	"testing"
)

func Test_Factory_Interface_Default(t *testing.T) {
	var _ Interface = Default()
}

func Test_Factory_Interface_Fake(t *testing.T) {
	var _ Interface = Fake()
}
