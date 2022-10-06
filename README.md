## Directory Structure

MVC Strucutre is being followed.

`main.go`: Driver Code for all of backend.

`models`: Database Models

`controllers`: Various handlers for HTTP Requests

`routes`: Definitions of Sub-routers

`utils`: Middlewares and other files common to multiple functions

`scripts`: Extra scripts (unrelated to the usual flow of the webapp)

`docs`: Documentation

## Routes

- `/mentor/`
    - `/form` POST - Register as a mentor.
    - `/dashboard` GET - Get all the dashboard stats of the mentors.
    - `/` GET - Get all the registered mentors
- `/student`
  - `/form` POST - Register a student for the event.
  - `/dashboard` GET - Get dashboard stats of a particular student.
  - `/blog` POST - Submit blog link for final submissions
  - `/stats` GET - Get all stats of all the users.
- `/project`
  - `/register` POST - Register a Project
  - `/` GET - Return all registered Porjects ! Open Route !
  - `/details` POST - Return all details of a particular project.
  - `/update` PUT - Update details of the project.
  - `/project` GET - Get Stats of a particular project.
  - `/dashboard` GET - Get dashbard stats of all the projects.
- `/oauth` POST - Oauth using Github Token to get details for participants

## File Naming Convention

- Separate file for each Model.

- Controller functions to be grouped together by their routes.

- Each Subroute has a separate file.

In short, don't keep any surprises. Use groupings as per your discretion.

## Dependencies

```
gorilla/mux
jinzhu/gorm
dgrijalva/jwt-go
```

Also uses `golanglint-ci` for linting code.

## Set up

- Clone the repo.

- You can use the Makefile for automating the commands. Run `make help` for a list of commands.

- Currently, only two(subject to change) commands are supported -

* `make lint` - Run the lint checks
* `make build` - For building the codebase

- Run `go get` to install the dependencies

- Run `go run main.go` in the terminal

- Run `go build` to build into a single binary.

- Run `gofmt -s -w .` to lint all the files in one go.
