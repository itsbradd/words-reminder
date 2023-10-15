package passhashing

import (
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) HashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *Service) VerifyPassword(hashedPassword, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(pass))
}

func New() fx.Option {
	var module = fx.Module("passhashing",
		fx.Provide(
			NewService,
		),
	)

	return module
}
