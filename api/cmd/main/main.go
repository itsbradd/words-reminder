package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sonngocme/words-reminder-be/bootstrap"
	"github.com/sonngocme/words-reminder-be/db"
	"github.com/sonngocme/words-reminder-be/pkg/jwt"
	"github.com/sonngocme/words-reminder-be/pkg/passhashing"
	usermodule "github.com/sonngocme/words-reminder-be/pkg/user"
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
