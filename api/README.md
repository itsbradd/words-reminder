# Words Reminder API

## Tech stack

- [Go](https://go.dev/) (Language)
- [Fiber](https://docs.gofiber.io/) (Web framework)
- [Sqlc](https://docs.sqlc.dev) (Database queries generator)
- [Goose](https://pressly.github.io/goose/) (Database migration management)

## Setup local development
```sh
# Install Goose migration management
$ go install github.com/pressly/goose/v3/cmd/goose@latest

# Install project dependencies
$ go mod download

# Start the server
$ go run .
```