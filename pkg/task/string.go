package task

import (
	"encoding/json"
)

func FromString(s string) *Task {
	t := &Task{}

	err := json.Unmarshal([]byte(s), t)
	if err != nil {
		panic(err)
	}

	return t
}

func ToString(t *Task) string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	return string(b)
}
