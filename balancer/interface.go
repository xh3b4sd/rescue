package balancer

type Interface interface {
	// Dev assesses the deviation between the current and desired distribution.
	// The resulting map contains key/value pairs for workers that must revoke
	// ownership to the specified amount of tasks in order to regain balance
	// within the system.
	Dev(cur map[string]int, des map[string]int) map[string]int
	// Opt returns the optimal distribution of tasks shared among the given
	// workers. While there might be different implementations, the most simple
	// one may look like the following. Given workers a, b and c. And given 7
	// tasks to balance between these three workers. The returned balance could
	// look like shown below.
	//
	//     a 3
	//     b 2
	//     c 2
	//
	Opt(wrk []string, tas int) map[string]int
}
