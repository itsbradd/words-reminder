package pkg

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

type AppRouter interface {
	Setup(fiber.Router)
}

func ProvideRouters(routers ...any) fx.Option {
	return fx.Provide(
		lo.Map(routers, func(item any, index int) any {
			return AsRoute(item)
		})...,
	)
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(AppRouter)),
		fx.ResultTags(`group:"routers"`),
	)
}
