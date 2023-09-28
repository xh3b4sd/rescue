package balancer

func Default() Interface {
	var err error

	var bal Interface
	{
		c := Config{}

		bal, err = New(c)
		if err != nil {
			panic(err)
		}
	}

	return bal
}
