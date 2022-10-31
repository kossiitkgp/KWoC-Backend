# Building and Running

- Please use go 1.19 or check `go.mod` for the required version
- Default port is 8080. To change it, set environment variable `BACKEND_PORT` to desired port number
- Run `./build.sh`. If it doesn't run, make sure it is executable


> To view the program as doc, run : `godoc -http=:6060` and checkout at `http://localhost:6060/pkg/kwoc-backend/`

# Module structure

```
├── cmd
│   └── backend.go
│   └── ...
├── controllers
│   └── index.go
│   └── ...
├── server
│   ├── router.go
│   └── routes.go
│   └── ...
└── utils
    └── logger.go
    └── ...
```

`cmd` : contains the entrypoint of the backend (main package) \\
`controllers` : handler functions for the routes defined \\
`server` : contains the router logic and routes \\
`utils` : contains misc functions like logger

> For new packages, a few suggestions

- For middlewares, please create and use `middleware` directory
- If there are any css,html or other static files, use `static` directory
- Do not keep many functions in utils, if they can be grouped in a package, then do so.


# Libraries to refer

- gorilla/mux : [https://github.com/gorilla/mux](https://github.com/gorilla/mux). Used for routing

> Helpful libraries :

- go-orm/gorm : [https://github.com/go-gorm/gorm](https://github.com/go-gorm/gorm). Used for database modelling

