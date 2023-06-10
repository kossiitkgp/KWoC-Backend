# KWoC Backend
KWoC backend server written in Go.

## Table of Contents
- [Development](#development)
  - [Setting Up Locally](#setting-up-locally)
  - [Building](#building)
  - [File Naming Convention](#file-naming-convention)
- [Project Structure](#project-structure)
  - [File Structure](#file-structure)
  - [Endpoints](#endpoints)
  - [Environment Variables](#environment-variables)

## Development
### Setting Up Locally
- Install [Go](https://go.dev).
- Clone this repository.
- Run `go get` in the repository to download all the dependencies.
- Create a `.env` file to store all the [environment variables](#environment-variables).
- Set the `DEV` environment variable to `true`. This makes the server use a local sqlite3 database, `devDB.db`.
- Run `go run .` to start the server.
- Optionally install [SQLiteStudio](https://sqlitestudio.pl/) or any similar tool to help manage the local sqlite3 database file `devDB.db`.

### File Naming Convention
See also: [File Structure](#file-structure).

- Create a separate file for each model in `models/`.
- Group the controller functions by their routes.
- Create a separate file for each subroute.

### Building
- Clone the repo and `cd` into its directory.
- Run the following commands
  ```sh
  docker-compose build
  docker-compose up
  ```
**NOTE**: If you face the following issue with `docker-compose`.
> strconv.Atoi: parsing "": invalid syntax

This is because `docker-compose` is creating an arbitrary container not in _docker-compose.yml_ and terminating the starting process.</br>
**FIX**: Use `docker-compose-v1` instead of `docker-compose`.

## Project Structure
### File Structure
This project follows the MVC Structure.
- `main.go`: The main file that is executed.
- `models/`: Database models.
- `controllers/`: Various handlers for HTTP requests.
- `routes/`: Definitions of sub-routers.
- `utils/`: Middlewares and other common utility functions.
- `scripts/`: Extra scripts (unrelated to the usual flow of the web app).
- `docs/`: Documentation.

### Endpoints
The API exposes the following endpoints. (See also: [File Structure](#file-structure))

#### Healthcheck
Files: `routes/healthcheck.go`, `controllers/healthcheck.go`.
- `/healthcheck`(GET): Responds with the server and database status.
- `/healthcheck/ping`(GET): Responds with "Pong!" if the request was successful.

#### Mentor
Files: `routes/mentor.go`, `controllers/mentor.go`.
- `/mentor/form`(POST): Registers a mentor. The following parameters are required.
  - `name`: The name of the mentor.
  - `email`: The email address of the mentor.
  - `username`: The Github username of the mentor.
- `/mentor/dashboard`(POST): Responds with the information to be displayed on the mentor's dashboard, including a list of projects. The following parameters are required.
  - `username`: The Github username of the mentor.
- `/mentor/all`(POST): Responds with a list of all the mentors' names and usernames.

#### OAuth
Files: `routes/oauth.go`, `controllers/UserOAuth.go`.
- `/oauth`(POST): Logs in a user via Github OAuth.

#### Project
Files: `routes/project.go`, `controllers/project.go`.
- `/project`(GET): Responds with a list of all projects.
- `/project/add`(POST): Adds a new project to the database. See the file `controllers/project.go` for a list of parameters.
- `/project/details`(POST): Responds with the details of a project. The following parameters are required.
  - `id`: The ID of the project in the database.
- `/project/update`(PUT): Updates the details of a project. See the file `controllers/project.go` for a list of parameters.

#### Stats
Files: `routes/stats.go`, `controllers/stats-overall.go`, `controllers/stats-project.go`, `controllers/stats-student.go`.
- `/stats/overall`(GET): Responds with the overall statistics.
- `/stats/projects`(GET): Responds with a list of all projects' statistics.
- `/stats/students`(GET): Responds with a list of all students' statistics.
- `/stats/student/exists/{username}`(GET): Responds with "true" if statistics for the given student exist in the database and "false" otherwise.
- `/stats/student/{username}`(GET): Responds with the stats for the given student.
- `/stats/mentor/{username}`(GET): Responds with a list of the mentor's project stats.

#### Student
Files: `routes/student.go`, `controllers/student.go`.
- `/student/form`(POST): Registers a student. The following parameters are required.
  - `name`: The name of the student.
  - `email`: The email address of the student.
  - `college`: The institute the student studies at.
  - `username`: The Github username of the student.
- `/student/dashboard`(POST): Responds with the information to be displayed on the student's dashboard. The following parameters are required.
  - `username`: The Github username of the student.
- `/student/bloglink`(POST): Adds a link to the student's blog to the database. The following parameters are required.
  - `username`: The Github username of the student.
  - `bloglink`: A link to the student's blog.

### Environment Variables
Environment variables can be set using a .env file (see the `.env.template` file). The following variables are used.
- `DEV`: When set to `true`, uses a local sqlite3 database from a file `devDB.db`.
- `DATABASE_USERNAME`: The username used to log into the database. (Valid when `DEV` is not set to `true`)
- `DATABASE_PASSWORD`: The password used to log into the database. (valid when `DEV` is not set to `true`)
- `DATABASE_NAME`: The name of the database to log into. (Valid when `DEV` is not set to `true`)
- `DATABASE_HOST`: The host/url used to log into the database. (Valid when `DEV` is not set to `true`)
- `DATABASE_PORT`: The port used to log into the database. (Valid when `DEV` is not set to `true`)
- `client_id`: The client id used for Github OAuth.
- `client_secret`: The client secret used for Github OAuth.
- `JWT_SECRET_KEY`: The secret key used to create a JWT token.

****
> Please update this documentation if you make changes to the KWoC backend or any other part of KWoC which affects the backend. Future humans will praise you.