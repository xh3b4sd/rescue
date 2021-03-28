package metric

type Collector interface {
	Dec()
	Dur(o func() error) error
	Get() float64
	Inc()
	Set(i float64)
}
