package key

import "fmt"

func Queue(que string) string {
	return fmt.Sprintf("rescue.io:%s", que)
}
