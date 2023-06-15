# KWoC Backend v2.0
[WIP] KWoC backend server v2.0 (also) written in Go (but better).

## Table of Contents
- [Development](#development)
  - [Setting Up Locally](#setting-up-locally)
  - [Building](#building)
  - [File Naming Convention](#file-naming-convention)
- [Project Structure](#project-structure)
  - [Libraries Used](#libraries-used)
  - [File Structure](#file-structure)
  - [Endpoints](#endpoints)
  - [Command-Line Arguments](#command-line-arguments)
  - [Environment Variables](#environment-variables)
  - [Github OAuth](#github-oauth)

## Development
### Setting Up Locally
[WIP]

### Building
[WIP]
- Please use go 1.19 or check `go.mod` for the required version.
- Default port is 8080. To change it, set environment variable `BACKEND_PORT` to desired port number.
- Run `./build.sh`. If it doesn't run, make sure it is executable.
> To view the program as doc, run : `godoc -http=:6060` and checkout at `http://localhost:6060/pkg/kwoc-backend/`

### File Naming Convention
[WIP]

## Project Structure
[WIP]
### Libraries Used
- gorilla/mux : [https://github.com/gorilla/mux](https://github.com/gorilla/mux). Used for routing
- go-orm/gorm : [https://github.com/go-gorm/gorm](https://github.com/go-gorm/gorm). Used for database modelling

### File Structure
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

- `cmd` : Contains the entrypoint of the backend (main package).
- `controllers` : Handler functions for the routes defined.
- `server` : Contains the router logic and routes.
- `utils` : Contains misc functions like logger.

- For middlewares, please create and use `middleware` directory.
- If there are any css,html or other static files, use `static` directory.
- Do not keep many functions in utils, if they can be grouped in a package, then do so.

### Endpoints
### Command-Line Arguments
The following command-line arguments are accepted by `cmd/backend.go`. `--argument=value`, `--argument value`, `-argument=value`, and `-argument value` are all acceptable formats to pass a value to the command-line argument.
- `envFile`: A file to load environment variables from. (Default: `.env`)

### Environment Variables
### Github OAuth

