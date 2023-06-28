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
- Install [Go](https://go.dev/).
- Clone this repository.
- Run `go mod download` in the cloned repository to download all the dependencies.
- Create a `.env` file to store all the [environment variables](#environment-variables). You can use the `.env.template` file for this.
- Set the `DEV` environment variable to `true`.
- Optionally set up [Github OAuth](#github-oauth) to test the endpoints which require login. (See also: [Endpoints](#endpoints))
- Run `go run cmd/backend.go` to start the server.
- Optionally install [SQLiteStudio](https://sqlitestudio.pl/) or a similar tool to help manage the local database `devDB.db` (sqlite3).
- Optionally install [Postman](https://www.postman.com/) or a similar tool to test the API endpoints.

### Building
[WIP]
- Please use go 1.19 or check `go.mod` for the required version.
- Default port is 8080. To change it, set environment variable `BACKEND_PORT` to desired port number.
- Run `./build.sh`. If it doesn't run, make sure it is executable.
> To view the program as doc, run : `godoc -http=:6060` and checkout at `http://localhost:6060/pkg/kwoc-backend/`

### Setting up pre-commit hook (optional)
- Check if `golangci-lint` is installed, if not, install from [golangci-lint](https://golangci-lint.run/usage/install/)
- Run `git config core.hooksPath .githooks`, see [core.hooksPath](https://git-scm.com/docs/git-config#Documentation/git-config.txt-corehooksPath)

### File Naming Convention
See also [File Structure](#file-structure).

1. Test Files: Tests for a particular `file.go` are placed in a file named `file_test.go` and are placed in the same directory.
2. Model Files: Database models are placed in the `models` directory. Each file corresponds to one table in the database and the name of the file corresponds with the name of the database table.

## Project Structure
[WIP]
### Libraries Used
- [gorilla/mux](https://github.com/gorilla/mux): Used for routing.
- [gorm.io/gorm](https://gorm.io): Used for database modelling.
- [joho/godotenv](https://github.com/joho/godotenv): Used for loading environment variables from a `.env` file.
- [rs/zerolog](https://github.com/rs/zerolog): Used for logging information, errors, and warnings.

### File Structure
```
├── cmd
│   ├── backend.go
│   └── ...
├── controllers
│   ├── index.go
│   └── ...
├── server
│   ├── router.go
│   ├── routes.go
│   └── ...
├── models
│   ├── mentors.go
│   └── ...
├── utils
│   ├── database.go
│   └── ...
└── middleware
    ├── logger.go
    └── ...
```

- `cmd` : Contains the entrypoint of the backend (main package).
- `controllers` : Handler functions for the routes defined.
- `server` : Contains the router logic and routes.
- `models` : Contains KWoC database models.
- `utils` : Contains misc functions like database utils.
- `middleware`: Contains all middleware.

- If there are any css,html or other static files, use `static` directory.
- Do not keep many functions in utils, if they can be grouped in a package, then do so.

### Endpoints
### Command-Line Arguments
The following command-line arguments are accepted by `cmd/backend.go`. `--argument=value`, `--argument value`, `-argument=value`, and `-argument value` are all acceptable formats to pass a value to the command-line argument.
- `envFile`: A file to load environment variables from. (Default: `.env`)

### Environment Variables
Environment variables can be set using a `.env` (See [Command-Line Arguments](#command-line-arguments) to use a different file) file. The following variables are used. (See the `.env.template` file for an example)
- `DEV`: When set to `true`, uses a local sqlite3 database from a file `devDB.db` (or `DEV_DB_PATH` if set).
- `DEV_DB_PATH`: The path to a local sqlite3 database file to be used in development. (Valid when `DEV` is set to `true`) (NOTE: `testDB.db` is used for testing)
- `BACKEND_PORT`: The port on which the backend server runs. (Default: `8080`)
- `DATABASE_USERNAME`: The username used to log into the database. (Valid when `DEV` is not set to `true`)
- `DATABASE_PASSWORD`: The password used to log into the database. (valid when `DEV` is not set to `true`)
- `DATABASE_NAME`: The name of the database to log into. (Valid when `DEV` is not set to `true`)
- `DATABASE_HOST`: The host/url used to log into the database. (Valid when `DEV` is not set to `true`)
- `DATABASE_PORT`: The port used to log into the database. (Valid when `DEV` is not set to `true`)
- `GITHUB_OAUTH_CLIENT_ID`: The client id used for Github OAuth. (See [Github OAuth](#github-oauth))
- `GITHUB_OAUTH_CLIENT_SECRET`: The client secret used for Github OAuth. (See [Github OAuth](#github-oauth))
- `JWT_SECRET_KEY`: The secret key used to create a JWT token. (It can be a randomly generated string)
- `JWT_VALIDITY_TIME`: The amount of time (in hours) for which the generated JWT tokens should be valid.

### Github OAuth
KWoC uses Github [OAuth](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/differences-between-github-apps-and-oauth-apps#about-github-apps-and-oauth-apps) for log in authentication instead of passwords.

To set up the KWoC server, a Github OAuth application has to be created and the client id and secret has to be set in the [environment variables](#environment-variables).

- Follow [this](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/creating-an-oauth-app) documentation to create an OAuth app. In the production server, use the `koss-service` account to create the application.
- Set the Homepage URL to `https://kwoc.kossiitkgp.org` and the authorization callback URL to `https://kwoc.kossiitkgp.org/oauth/` in the production application.
- Copy the client ID and the client secret (this should NEVER be made public) and set the `GITHUB_OAUTH_CLIENT_ID` and `GITHUB_OAUTH_CLIENT_SECRET` [environment variables](#environment-variables) to these values.
