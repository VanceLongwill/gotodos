# Go Todos
> a simple todo list API written in Golang

## Features

- Read/Write capability over API
- Authenticated routes/requests
- Use database for storage
- JSON – format for communication
- Includes simple documentation (README)
- Includes unit tests

## Stack

- [Postgres](https://www.postgresql.org) for a database
- Golang's [`database/sql`](https://golang.org/pkg/database/sql/) and this [Postgres driver](https://github.com/lib/pq) for database interactions
- [Gin](https://github.com/gin-gonic/gin) for a http framework / routing
- [Docker](https://www.docker.com) for containerization
- Authentication using [JWT tokens](https://jwt.io)
- Golang dependency management with [dep](https://github.com/golang/dep)

## Installation

### Using the Makefile

#### Available commands

- `make` Builds the project locally
- `make test` Runs all tests
- `make lint` Lints the whole project using *gometalinter* (which is automatically installed if not already)
- `make coverage` Generates a HTML code coverage report [coverage_report.html](./coverage_report.html)
- `make run` runs main.go to connect to the db via localhost for development/debugging purposes

## Running locally

- Download the necessary dependencies with [dep](https://github.com/golang/dep)
```sh
dep ensure
```

- Set the relevant environment variables (see .env.sample)
```sh
cp .env.sample .env
source .env
```

- Launch the postgres image
```sh
docker-compose up --build db
```

- Launch the app, connecting to the db via localhost, by using 
```sh
make run
```

## Running in production

- Set the relevant environment variables (see .env.sample)
```sh
cp .env.sample .env
source .env
```

- Start all containers in background
```sh
docker-compose up --build
```

- Stop all running containers
```sh
docker-compose stop
```

- Remove all containers & volumes
```sh
docker-compose down
```

### Resources

#### Application Users

- **POST** `/api/v1/user/register` Register a user and obtain a JWT token

  Example

  ```sh
  curl -X POST localhost:8080/api/v1/user/register \
  --data '{"email":"test@google.com","firstName":"test","lastName":"user","password":"testpassword"}'
  ```

  Response

  ```json
  {
    "email": "test@google.com",
    "message": "User registered successfully!",
    "resourceId": 3,
    "status": 201,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MywiZXhwIjoxNTU0MzEzODk3fQ.iCqbCYPPJz6HHN_-IR2xmiGYuDeeyMDmxIFobaYUJA4"
  }
  ```
- **POST** `/api/v1/user/login` Login a user and obtain a JWT token

    > NB: JWT is set in cookies on compatible clients (e.g. browsers, postman)

  Example

  ```sh
  curl -X POST localhost:8080/api/v1/user/login \
  --data '{"email":"test@google.com","password":"testpassword"}'
  ```

  Response

  ```json
  {
    "email": "test@google.com",
    "message": "User logged in successfully!",
    "resourceId": 3,
    "status": 200,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MywiZXhwIjoxNTU0MzEzOTk2fQ.990Y5wV9YkMF-BDzFRQLGZL0zw0mgC820l_ZiDGdjFc"
  }
  ```


#### Todos (requires authentication)

- **GET** `/api/v1/todos/` Retrieves a list of a user's todos

  Example
  ```sh
  curl localhost:8080/api/v1/todos/ \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNTUzODg4OTkxfQ.eNn7bMfqHwA1ZF8Q87Ut0kdyZPntURuIGNuHMTvefJ8"
  ```
  Response
  ```json
  {
    "data": [
      {
        "createdAt": "2019-03-15T10:13:31.996901Z",
        "dueAt": "2019-03-15T10:13:31.996901Z",
        "id": 2,
        "isDone": false,
        "modifiedAt": "2019-03-15T10:13:31.996901Z",
        "note": "Some note",
        "title": "Some title"
      },
      {
        "createdAt": "2019-03-19T15:49:29.693415Z",
        "id": 9,
        "isDone": false,
        "modifiedAt": "2019-03-19T15:49:29.693415Z",
        "note": "Some note",
        "title": "Some title"
      },
    ],
    "status": 200
  }
  ```
- **POST** `/api/v1/todos/` Creates a new todo item for a user

  Example
  ```sh
  curl -X POST localhost:8080/api/v1/todos/ \
  --data '{"title":"some title","note":"some note"}' \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNTUzODg4OTkxfQ.eNn7bMfqHwA1ZF8Q87Ut0kdyZPntURuIGNuHMTvefJ8"
  ```
  Response (assuming the token is valid)
  ```json
  {
    "message": "Todo item created successfully!",
    "resourceId": 0,
    "status": 201
  }
  ```
- **GET** `/api/v1/todos/:id` Retrieves a single todo item by id

  Example
  ```sh
  curl localhost:8080/api/v1/todos/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNTUzODg4OTkxfQ.eNn7bMfqHwA1ZF8Q87Ut0kdyZPntURuIGNuHMTvefJ8"
  ```
  Response
  ```json
  {
    "data": {
      "createdAt": "2019-03-15T10:13:28.123377Z",
      "dueAt": "2019-03-15T10:13:28.123377Z",
      "id": 1,
      "isDone": false,
      "modifiedAt": "2019-03-15T10:13:28.123377Z",
      "note": "Some note",
      "title": "Some title"
    },
    "status": 200
  }
  ```
- **DELETE** `/api/v1/todos/:id` Deletes a single todo item by id

  Example
  ```sh
  curl -X DELETE localhost:8080/api/v1/todos/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNTUzODg4OTkxfQ.eNn7bMfqHwA1ZF8Q87Ut0kdyZPntURuIGNuHMTvefJ8"
  ```
  Response 
  ```json
  {
    "message": "Todo deleted successfully!",
    "resourceId": 1,
    "status": 200
  }
  ```
- **PUT** `/api/v1/todos/:id` Updates the title and note for a single todo item

  Example
  ```sh
  curl -X PUT localhost:8080/api/v1/todos/5 \
  --data '{"title":"changed the title","note":"changed the note"}' \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNTUzODg4OTkxfQ.eNn7bMfqHwA1ZF8Q87Ut0kdyZPntURuIGNuHMTvefJ8"
  ```

- **GET** `api/v1/todos/:id/completed` Marks a single todo item as completed

  Example
  ```sh
  curl localhost:8080/api/v1/todos/5/completed \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNTUzODg4OTkxfQ.eNn7bMfqHwA1ZF8Q87Ut0kdyZPntURuIGNuHMTvefJ8"
  ```

  Response
  ```json
  {
    "message": "Todo marked as complete successfully",
    "status": 200
  }
  ```

## DB Admin

- open a `psql` shell in the container (this doesn't require a local postgres installation)
  ```sh
  docker exec --tty --interactive \
  gotodos_db_1 psql \
  -h localhost -U gotodos -d postgres
  ```

***

### Licence 

MIT
