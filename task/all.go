package task

func (t *Task) All(key ...string) *Task {
	var tas *Task
	{
		tas = &Task{}
	}

	for _, x := range key {
		m, e := t.has(x)
		if !e {
			return nil
		}

		if tas.Meta == nil {
			tas.Meta = map[string]string{}
		}

		for k, v := range m {
			tas.Meta[k] = v
		}
	}

	return tas
}
