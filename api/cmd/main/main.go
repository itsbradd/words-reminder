package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itsbradd/words-reminder-be/bootstrap"
	"github.com/itsbradd/words-reminder-be/db"
	"github.com/itsbradd/words-reminder-be/pkg/jwt"
	"github.com/itsbradd/words-reminder-be/pkg/passhashing"
	usermodule "github.com/itsbradd/words-reminder-be/pkg/user"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fx.New(
		bootstrap.ProvideHTTPServer(),
		bootstrap.ProvideDBConn(),
		usermodule.New(),
		passhashing.New(),
		jwt.New(),
		fx.Invoke(func(*fiber.App, *db.Queries) {}),
	).Run()
}
