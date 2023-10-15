package user

import (
	"github.com/sonngocme/words-reminder-be/pkg"
	"github.com/sonngocme/words-reminder-be/pkg/passhashing"
	"go.uber.org/fx"
)

func New() fx.Option {
	var module = fx.Module("user",
		pkg.ProvideRouters(NewRouter),
		fx.Provide(
			NewHandler,
			fx.Annotate(
				NewService,
				fx.From(new(Storage), new(*passhashing.Service)),
			),
			NewStorage,
		),
	)

	return module
}
