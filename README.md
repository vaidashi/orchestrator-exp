
This repo was for experimenting with building an API in Go and constructing an orchestrator using the Docker API. 

#### Create a .env and populate it with: 


```
WORKER_HOST=localhost
WORKER_PORT=5555
MANAGER_HOST=localhost
MANAGER_PORT=5556
```

#### To run

```go run main.go```

#### To submit tasks:

```curl -v -X POST localhost:5556/tasks -d @task.json5556/tasks -d @task.json```

#### To query manager to assess status of tasks:

```curl localhost:5556/tasks```

