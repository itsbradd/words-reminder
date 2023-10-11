package bootstrap

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/sonngocme/words-reminder-be/pkg"
	"go.uber.org/fx"
	"log"
	"os"
)

func NewHTTPServer(lc fx.Lifecycle, routers []pkg.AppRouter) *fiber.App {
	app := fiber.New()
	api := app.Group("/api")
	v1 := api.Group("/v1")

	for _, router := range routers {
		router.Setup(v1)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				port := os.Getenv("PORT")
				if err := app.Listen(port); err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})
	return app
}

func ProvideHTTPServer() fx.Option {
	return fx.Provide(
		fx.Annotate(
			NewHTTPServer,
			fx.ParamTags(``, `group:"routers"`),
		),
	)
}
