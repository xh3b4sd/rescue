package task

func (t *Task) Any(key ...string) *Task {
	var tas *Task
	{
		tas = &Task{}
	}

	for _, x := range key {
		m, e := t.has(x)
		if !e {
			continue
		}

		if tas.Meta == nil {
			tas.Meta = map[string]string{}
		}

		for k, v := range m {
			tas.Meta[k] = v
		}
	}

	if len(tas.Meta) == 0 {
		return nil
	}

	return tas
}
