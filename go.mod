module github.com/xh3b4sd/rescue

go 1.22.5

require (
	github.com/google/go-cmp v0.7.0
	github.com/google/uuid v1.6.0
	github.com/prometheus/client_golang v1.22.0
	github.com/xh3b4sd/breakr v0.1.0
	github.com/xh3b4sd/logger v0.9.0
	github.com/xh3b4sd/objectid v0.2.0
	github.com/xh3b4sd/redigo v0.41.0
	github.com/xh3b4sd/tracer v0.11.1
)

require (
	github.com/FZambia/sentinel v1.1.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-redsync/redsync/v4 v4.13.0 // indirect
	github.com/gomodule/redigo v1.9.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	golang.org/x/sys v0.30.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)

retract [v0.0.0, v0.14.0]
