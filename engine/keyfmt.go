package engine

import "fmt"

func (e *Engine) Keyfmt() string {
	return fmt.Sprintf("rescue.io%s%s", e.sep, e.que)
}
