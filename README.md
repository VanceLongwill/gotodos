# Go Todos
> a simple todo list API written in Golang

## Features

- Read/Write capability over API
- Simple authentication of requests
- Use database for storage
- JSON – format for communication
- Includes simple documentation (Readme.md)
- Include Unit tests

## Installation

## DB

- pull the postgres docker images using 
    ```
    docker pull postgres
    ```

- create a volume to persist data
    ```
    docker volume create gotodos
    ```

- run the postgres container
  ```
  docker run -d --rm \
  --name pg_gotodos \
  -e POSTGRES_PASSWORD=gotodos \
  -p 5000:5000 \
  --mount source=gotodos,target=/var/lib/postgresql/data \
  postgres
  ```

- open a `psql` shell in the container (this doesn't require a local postgres installation)
  ```
  docker exec --tty --interactive \
  pg_gotodos psql \
  -h localhost -U postgres -d postgres
  ```

***

### Licence 

MIT
