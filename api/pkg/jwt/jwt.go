package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"maps"
	"os"
	"time"
)

type MapClaims jwt.MapClaims

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type GenClaimOpts func(claims MapClaims)

func (s *Service) NewWithClaims(claims MapClaims, opts ...GenClaimOpts) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	copyClaims := s.GenClaimOptions(claims, opts...)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(copyClaims))
	return token.SignedString([]byte(secret))
}

func (s *Service) GenClaimOptions(claims MapClaims, opts ...GenClaimOpts) MapClaims {
	copyClaims := maps.Clone(claims)
	for _, optFunc := range opts {
		optFunc(copyClaims)
	}
	return copyClaims
}

func (s *Service) GenIssuerClaim(val string) GenClaimOpts {
	return func(claims MapClaims) {
		claims["iss"] = val
	}
}

func (s *Service) GenSubjectClaim(val any) GenClaimOpts {
	return func(claims MapClaims) {
		claims["sub"] = val
	}
}

func (s *Service) GenAudienceClaim(val string) GenClaimOpts {
	return func(claims MapClaims) {
		claims["aud"] = val
	}
}

func (s *Service) GenExpireTimeClaim(val time.Time) GenClaimOpts {
	return func(claims MapClaims) {
		claims["exp"] = val.UnixMilli()
	}
}

func (s *Service) GenIssueAtClaim(val time.Time) GenClaimOpts {
	return func(claims MapClaims) {
		claims["iat"] = val.UnixMilli()
	}
}

func New() fx.Option {
	var module = fx.Module("jwt",
		fx.Provide(
			NewService,
		),
	)

	return module
}
