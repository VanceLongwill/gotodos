# Go Todos
> a simple todo list API written in Golang

## Features

- Read/Write capability over API
- Simple authentication of requests
- Use database for storage
- JSON – format for communication
- Includes simple documentation (Readme.md)
- Include Unit tests

## Stack

- Postgres for a database
- Gorm for an database ORM
- Gin for a http framework / routing
- Docker for containerization

## Installation

## Running Locally

- start all containers in background
```
docker-compose up -d
```

- stop all running containers
```
docker-compose stop
```

- remove all containers & volumes
```
docker-compose down
```

### Endpoints

- **GET** `/api/v1/todos/`

  Example
  ```
  curl localhost:8080/api/v1/todos/
  ```
- **POST** `/api/v1/todos/`

  Example
  ```
  curl -X POST localhost:8080/api/v1/todos/-d title="second" -d note="hello world"
  ```
- **GET** `/api/v1/todos/:id`

  Example
  ```
  curl localhost:8080/api/v1/todos/
  ```
- **DELETE** `/api/v1/todos/:id`

  Example
  ```
  curl -X DELETE localhost:8080/api/v1/todos/3
  ```
- **PUT** `/api/v1/todos/:id`

  Example
  ```
  curl -X PUT localhost:8080/api/v1/todos/1 -d title="second" -d note="hello world"
  ```

## DB Admin

- open a `psql` shell in the container (this doesn't require a local postgres installation)
  ```
  docker exec --tty --interactive \
  gotodos_db_1 psql \
  -h localhost -U gotodos -d postgres
  ```

***

### Licence 

MIT
