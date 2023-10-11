package user

import (
	"github.com/sonngocme/words-reminder-be/pkg"
	"go.uber.org/fx"
)

func New() fx.Option {
	var module = fx.Module("user",
		pkg.ProvideRouters(NewRouter),
		fx.Provide(
			NewHandler,
			fx.Annotate(
				NewService,
				fx.As(new(service)),
			),
			NewStorage,
		),
	)

	return module
}
