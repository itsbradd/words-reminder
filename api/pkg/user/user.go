package user

import (
	"github.com/sonngocme/words-reminder-be/pkg"
	"github.com/sonngocme/words-reminder-be/pkg/jwt"
	"github.com/sonngocme/words-reminder-be/pkg/passhashing"
	"go.uber.org/fx"
)

func New() fx.Option {
	var module = fx.Module("user",
		pkg.ProvideRouters(NewRouter),
		fx.Provide(
			fx.Annotate(
				NewHandler,
				fx.From(new(Service), new(*jwt.Service)),
			),
			fx.Annotate(
				NewService,
				fx.From(new(Storage), new(*passhashing.Service)),
			),
			NewStorage,
		),
	)

	return module
}
