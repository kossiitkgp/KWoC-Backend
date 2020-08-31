### Logging
For the case of error handling, a logger with name `**LOG**` has been defined in `utils/logging.go`, with configuration set to logging at STDERR level currently.


To log, import the `**LOG**` from utils package. And use the Println method to log the error


### Makefile

#### Build the codebase
```make build```

###### Commands executed
```go build```

#### Linting
```make lint```

###### Commands executed
```go fmt -s -w .```
