package rescue

type Interface interface {
	Claim() (Task, error)
	Purge(t Task) error
}
