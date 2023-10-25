package user

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/sonngocme/words-reminder-be/db"
	"github.com/sonngocme/words-reminder-be/pkg"
	"github.com/sonngocme/words-reminder-be/pkg/jwt"
	"time"
)

type Service interface {
	CreateUser(context.Context, db.CreateUserParams) (int64, error)
	HashPassword(string) (string, error)
	SetUserRefreshToken(ctx context.Context, id int64, token string) error
	Login(ctx context.Context, info LoginInfo) (*Credentials, error)
	SignUp(ctx context.Context, info SignUpInfo) (*Credentials, error)
}

type JWTService interface {
	NewWithClaims(jwt.MapClaims, ...jwt.GenClaimOpts) (string, error)
	GenIssuerClaim(val string) jwt.GenClaimOpts
	GenSubjectClaim(val any) jwt.GenClaimOpts
	GenAudienceClaim(val string) jwt.GenClaimOpts
	GenIssueAtClaim(val time.Time) jwt.GenClaimOpts
	GenExpireTimeClaim(val time.Time) jwt.GenClaimOpts
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
		return pkg.ErrParseReqBody
	}

	err := signUpInfo.Validate()
	if err != nil {
		return pkg.NewBodyValidationErr(err.(validation.Errors))
	}

	credentials, err := h.s.SignUp(c.Context(), *signUpInfo)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(pkg.SuccessRes[Credentials]{
		Message: "Signup success!",
		Data:    *credentials,
	})
}

func (h Handler) Login(c *fiber.Ctx) error {
	loginInfo := new(LoginInfo)
	if err := c.BodyParser(loginInfo); err != nil {
		return pkg.ErrParseReqBody
	}
	err := loginInfo.Validate()
	if err != nil {
		return pkg.NewBodyValidationErr(err.(validation.Errors))
	}

	credentials, err := h.s.Login(c.Context(), *loginInfo)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(pkg.SuccessRes[Credentials]{
		Message: "login success",
		Data:    *credentials,
	})
}
