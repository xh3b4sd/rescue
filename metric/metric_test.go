package metric

import "testing"

func Test_Metric(t *testing.T) {
	m := Default()

	{
		i := m.Engine.Create.Cal.Get()
		if i != 0 {
			t.Fatal("i must be 0")
		}
	}

	{
		i := m.Engine.Delete.Cal.Get()
		if i != 0 {
			t.Fatal("i must be 0")
		}
	}

	{
		m.Engine.Create.Cal.Inc()
	}

	{
		i := m.Engine.Create.Cal.Get()
		if i != 1 {
			t.Fatal("i must be 1")
		}
	}

	{
		i := m.Engine.Expire.Cal.Get()
		if i != 0 {
			t.Fatal("i must be 0")
		}
	}

	{
		m.Engine.Create.Cal.Set(5)
	}

	{
		i := m.Engine.Create.Cal.Get()
		if i != 5 {
			t.Fatal("i must be 5")
		}
	}

	{
		m.Engine.Create.Cal.Dec()
		m.Engine.Create.Cal.Dec()
		m.Engine.Create.Cal.Dec()
	}

	{
		i := m.Engine.Create.Cal.Get()
		if i != 2 {
			t.Fatal("i must be 2")
		}
	}

	{
		i := m.Engine.Search.Cal.Get()
		if i != 0 {
			t.Fatal("i must be 0")
		}
	}

	{
		m.Reset()
	}

	{
		i := m.Engine.Create.Cal.Get()
		if i != 0 {
			t.Fatal("i must be 0")
		}
	}
}
