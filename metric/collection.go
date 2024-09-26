package metric

type Collection struct {
	Engine *CollectionEngine
	Task   *CollectionTask
}

type CollectionEngine struct {
	Create *CollectionEngineCollector
	Cycles *CollectionEngineCollector
	Delete *CollectionEngineCollector
	Exists *CollectionEngineCollector
	Expire *CollectionEngineCollector
	Extend *CollectionEngineCollector
	Lister *CollectionEngineCollector
	Search *CollectionEngineCollector
	Ticker *CollectionEngineCollector
}

type CollectionEngineCollector struct {
	Cal Interface
	Dur Interface
	Err Interface
}

type CollectionTask struct {
	Expired  Interface
	Extended Interface
	Inactive Interface
	NotFound Interface
	Obsolete Interface
	Outdated Interface
	Parallel Interface
}

func (c *Collection) Reset() {
	c.Engine.Create.Cal.Res()
	c.Engine.Create.Dur.Res()
	c.Engine.Create.Err.Res()

	c.Engine.Cycles.Cal.Res()
	c.Engine.Cycles.Dur.Res()
	c.Engine.Cycles.Err.Res()

	c.Engine.Delete.Cal.Res()
	c.Engine.Delete.Dur.Res()
	c.Engine.Delete.Err.Res()

	c.Engine.Exists.Cal.Res()
	c.Engine.Exists.Dur.Res()
	c.Engine.Exists.Err.Res()

	c.Engine.Expire.Cal.Res()
	c.Engine.Expire.Dur.Res()
	c.Engine.Expire.Err.Res()

	c.Engine.Extend.Cal.Res()
	c.Engine.Extend.Dur.Res()
	c.Engine.Extend.Err.Res()

	c.Engine.Lister.Cal.Res()
	c.Engine.Lister.Dur.Res()
	c.Engine.Lister.Err.Res()

	c.Engine.Search.Cal.Res()
	c.Engine.Search.Dur.Res()
	c.Engine.Search.Err.Res()

	c.Engine.Ticker.Cal.Res()
	c.Engine.Ticker.Dur.Res()
	c.Engine.Ticker.Err.Res()

	c.Task.Expired.Res()
	c.Task.Extended.Res()
	c.Task.Inactive.Res()
	c.Task.NotFound.Res()
	c.Task.Obsolete.Res()
	c.Task.Outdated.Res()
	c.Task.Parallel.Res()
}
