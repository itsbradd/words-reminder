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

# Install Sqlc database queries generator
$ go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Install project dependencies
$ go mod download

# Start the server
$ go run .
```

## Migrations management (Local)
```sh
# Create a new migration file
$ goose -dir sql/migrations create {migration_name} sql

# Apply migrations
$ goose -dir sql/migrations mysql "root:thisisverysecret@/words_reminder?parseTime=true" up

# Apply one next migration
$ goose -dir sql/migrations mysql "root:thisisverysecret@/words_reminder?parseTime=true" up-by-one

# Revert to previous migration
$ goose -dir sql/migrations mysql "root:thisisverysecret@/words_reminder?parseTime=true" redo

# Show current migration version
$ goose -dir sql/migrations mysql "root:thisisverysecret@/words_reminder?parseTime=true" version

# Show current migration status
$ goose -dir sql/migrations mysql "root:thisisverysecret@/words_reminder?parseTime=true" status
```

## Generate database queries
```sh
$ sqlc generate
```