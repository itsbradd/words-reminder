package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sonngocme/words-reminder-be/bootstrap"
	usermodule "github.com/sonngocme/words-reminder-be/pkg/user"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		bootstrap.ProvideHTTPServer(),
		bootstrap.ProvideDBConn(),
		usermodule.New(),
		fx.Invoke(func(*fiber.App) {}),
	).Run()
}
