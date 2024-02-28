module github.com/xh3b4sd/rescue

go 1.21

require (
	github.com/google/go-cmp v0.6.0
	github.com/google/uuid v1.6.0
	github.com/prometheus/client_golang v1.19.0
	github.com/xh3b4sd/breakr v0.1.0
	github.com/xh3b4sd/logger v0.8.1
	github.com/xh3b4sd/redigo v0.37.0
	github.com/xh3b4sd/tracer v0.11.1
)

require (
	github.com/FZambia/sentinel v1.1.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/go-redsync/redsync/v4 v4.11.0 // indirect
	github.com/gomodule/redigo v1.8.9 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

retract [v0.0.0, v0.14.0]
