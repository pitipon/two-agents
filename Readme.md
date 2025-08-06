# Run Redis
```
docker run -p 6379:6379 redis
```

# Run both agent
* open terminal 1 
```
cd AgentA
go run main.go
```

* open terminal 2
```
cd AgentB
go run main.go
```