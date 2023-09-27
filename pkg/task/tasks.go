package task

import "strings"

type Tasks []*Task

// With returns a list of tasks matching any metadata identified by the list of
// provided prefixes. If metadata does not exist for a prefix, nil may be
// returned. That means that the returned list of tasks will be nil, unless any
// provided prefix could match.
func (t Tasks) With(pre ...string) []*Task {
	var tas []*Task
	{
		tas = []*Task{}
	}

	for _, x := range t {
		for k := range x.Obj.Metadata {
			if prefix(pre, k) {
				tas = append(tas, x)
				break
			}
		}
	}

	if len(tas) == 0 {
		return nil
	}

	return tas
}

func prefix(pre []string, key string) bool {
	for _, p := range pre {
		if strings.HasPrefix(key, p) {
			return true
		}
	}

	return false
}
