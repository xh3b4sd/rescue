package task

func (t *Task) IncBackoff(d int) {
	b := t.GetBackoff()

	b += d

	t.SetBackoff(b)
}

func (t *Task) IncVersion(d int) {
	v := t.GetVersion()

	v += d

	t.SetVersion(v)
}
