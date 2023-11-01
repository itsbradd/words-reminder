package user

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/itsbradd/words-reminder-be/db"
	"github.com/itsbradd/words-reminder-be/pkg"
)

type Service interface {
	CreateUser(context.Context, db.CreateUserParams) (int64, error)
	SetUserRefreshToken(ctx context.Context, id int64, token string) error
	Login(ctx context.Context, info LoginInfo) (*Credentials, error)
	SignUp(ctx context.Context, info SignUpInfo) (*Credentials, error)
	RefreshAccessToken(ctx context.Context, info RefreshAccessTokenInfo) (*AccessToken, error)
}

type Handler struct {
	s Service
}

func NewHandler(s Service) Handler {
	return Handler{
		s: s,
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

	return c.Status(fiber.StatusOK).JSON(pkg.Response[Credentials]{
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

	return c.Status(fiber.StatusOK).JSON(pkg.Response[Credentials]{
		Message: "login success",
		Data:    *credentials,
	})
}

func (h Handler) RefreshAccessToken(c *fiber.Ctx) error {
	refreshAccessTokenInfo := new(RefreshAccessTokenInfo)
	if err := c.BodyParser(refreshAccessTokenInfo); err != nil {
		return pkg.ErrParseReqBody
	}
	err := refreshAccessTokenInfo.Validate()
	if err != nil {
		return pkg.NewBodyValidationErr(err.(validation.Errors))
	}

	accessToken, err := h.s.RefreshAccessToken(c.Context(), *refreshAccessTokenInfo)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(pkg.Response[AccessToken]{
		Message: "refresh access token succeed",
		Data:    *accessToken,
	})
}
