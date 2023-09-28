# rescue

Reconciliation driven resource queue.



### Conformance Tests

```
docker run --rm --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
```

```
go test ./... -race -tags redis
```
