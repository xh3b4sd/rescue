package metric

import "testing"

func Test_Metric(t *testing.T) {
	m := New()

	{
		i := m.Engine.Create.Action.Get()
		if i != 0 {
			t.Fatal("i must be 0")
		}
	}

	{
		i := m.Engine.Delete.Action.Get()
		if i != 0 {
			t.Fatal("i must be 0")
		}
	}

	{
		m.Engine.Create.Action.Inc()
	}

	{
		i := m.Engine.Create.Action.Get()
		if i != 1 {
			t.Fatal("i must be 1")
		}
	}

	{
		i := m.Engine.Expire.Action.Get()
		if i != 0 {
			t.Fatal("i must be 0")
		}
	}

	{
		m.Engine.Create.Action.Set(5)
	}

	{
		i := m.Engine.Create.Action.Get()
		if i != 5 {
			t.Fatal("i must be 5")
		}
	}

	{
		i := m.Engine.Metric.Action.Get()
		if i != 0 {
			t.Fatal("i must be 0")
		}
	}

	{
		m.Engine.Create.Action.Dec()
		m.Engine.Create.Action.Dec()
		m.Engine.Create.Action.Dec()
	}

	{
		i := m.Engine.Create.Action.Get()
		if i != 2 {
			t.Fatal("i must be 2")
		}
	}

	{
		i := m.Engine.Search.Action.Get()
		if i != 0 {
			t.Fatal("i must be 0")
		}
	}

	{
		m.Reset()
	}

	{
		i := m.Engine.Create.Action.Get()
		if i != 0 {
			t.Fatal("i must be 0")
		}
	}
}
