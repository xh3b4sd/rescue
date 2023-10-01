package key

import "fmt"

// TODO make this part of the engine interface Engine.Keyfmt
func Queue(que string) string {
	return fmt.Sprintf("rescue.io:%s", que)
}
