package task

func (t *Task) IncBackoff(d int) {
	b := t.GetBackoff()

	b += d

	t.SetBackoff(b)
}

func (t *Task) IncExpire(d int64) {
	e := t.GetExpire()

	e += d

	t.SetExpire(e)
}

func (t *Task) IncVersion(d int) {
	v := t.GetVersion()

	v += d

	t.SetVersion(v)
}
