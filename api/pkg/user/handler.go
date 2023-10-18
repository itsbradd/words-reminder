package user

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/sonngocme/words-reminder-be/db"
	"github.com/sonngocme/words-reminder-be/pkg/jwt"
	"time"
)

type Service interface {
	SignUpUser(context.Context, db.SignUpUserParams) (int64, error)
	HashPassword(string) (string, error)
	SetUserRefreshToken(ctx context.Context, id int64, token string) error
}

type JWTService interface {
	NewWithClaims(jwt.MapClaims, ...jwt.GenClaimOpts) (string, error)
	GenIssuerClaim(val string) jwt.GenClaimOpts
	GenSubjectClaim(val any) jwt.GenClaimOpts
	GenAudienceClaim(val string) jwt.GenClaimOpts
	GenIssueAtClaim(val time.Time) jwt.GenClaimOpts
}

type Handler struct {
	s   Service
	jwt JWTService
}

func NewHandler(s Service, jwt JWTService) Handler {
	return Handler{
		s:   s,
		jwt: jwt,
	}
}

func (h Handler) SignUp(c *fiber.Ctx) error {
	signUpInfo := new(SignUpInfo)
	if err := c.BodyParser(signUpInfo); err != nil {
		return err
	}

	err := signUpInfo.Validate()
	if err != nil {
		return err
	}

	hashedPass, err := h.s.HashPassword((*signUpInfo).Password)
	if err != nil {
		return err
	}

	userId, err := h.s.SignUpUser(context.Background(), db.SignUpUserParams{
		Username: (*signUpInfo).Username,
		Password: hashedPass,
	})
	if err != nil {
		return err
	}

	token, err := h.jwt.NewWithClaims(jwt.MapClaims{},
		h.jwt.GenIssuerClaim("Words Reminder"),
		h.jwt.GenSubjectClaim(userId),
		h.jwt.GenAudienceClaim("Words Reminder"),
		h.jwt.GenIssueAtClaim(time.Now()),
	)
	if err != nil {
		return err
	}

	err = h.s.SetUserRefreshToken(context.Background(), userId, token)
	if err != nil {
		return nil
	}
	return c.SendString(token)
}
