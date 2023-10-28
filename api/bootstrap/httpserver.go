package bootstrap

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/itsbradd/words-reminder-be/pkg"
	"go.uber.org/fx"
	"log"
	"os"
)

func NewHTTPServer(lc fx.Lifecycle, routers []pkg.AppRouter) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			var e *fiber.Error
			type FailResponse interface {
				GetStatusCode() int
				GetErrorCode() int
				Error() string
			}
			var failRes FailResponse
			if errors.As(err, &failRes) {
				return ctx.Status(failRes.GetStatusCode()).JSON(failRes)
			} else if errors.As(err, &e) {
				return ctx.Status(e.Code).SendString(e.Error())
			}

			return nil
		},
	})
	app.Get("/swagger/*", swagger.HandlerDefault)

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
