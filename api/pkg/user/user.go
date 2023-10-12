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
			NewService,
			NewStorage,
		),
	)

	return module
}
