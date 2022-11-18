## Directory Structure

MVC Strucutre is being followed.

`main.go`: Driver Code for all of backend.

`models`: Database Models

`controllers`: Various handlers for HTTP Requests

`routes`: Definitions of Sub-routers

`utils`: Middlewares and other files common to multiple functions

`scripts`: Extra scripts (unrelated to the usual flow of the webapp)

`docs`: Documentation

## File Naming Convention

- Separate file for each Model.

- Controller functions to be grouped together by their routes.

- Each Subroute has a separate file.

In short, don't keep any surprises. Use groupings as per your discretion.

## Dependencies

```
gorilla/mux
jinzhu/gorm
golang-jwt/jwt/v4
```

Also uses `golanglint-ci` for linting code.

## Building

- Clone the repo and `cd` into it's directory.
- Run the following commands
  ```sh
  docker-compose build
  docker-compose up 
  ```
**NOTE**: If you face the following issue with `docker-compose`. 
> strconv.Atoi: parsing "": invalid syntax

This is because `docker-compose` is creating arbitrary container not in _docker-compose.yml_ and terminating the starting process.</br>
**FIX**: Use `docker-compose-v1` instead of `docker-compose`.

<!-- These are steps before the dockerfile was created, basically without it

- You can use the Makefile for automating the commands. Run `make help` for a list of commands.

- Currently, only two(subject to change) commands are supported - 
* `make lint` - Run the lint checks
* `make build` - For building the codebase

- Run `go get` to install the dependencies

- Run `go run main.go` in the terminal

- Run `go build` to build into a single binary.

- Run `gofmt -s -w .` to lint all the files in one go.
-->
