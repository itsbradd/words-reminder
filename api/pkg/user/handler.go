package user

import (
	"context"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/sonngocme/words-reminder-be/db"
)

type service interface {
	SignUpUser(context.Context, db.SignUpUserParams) (sql.Result, error)
}

type Handler struct {
	s service
}

func NewHandler(s service) Handler {
	return Handler{
		s: s,
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

	_, err = h.s.SignUpUser(context.Background(), db.SignUpUserParams{
		Username: (*signUpInfo).Username,
		Password: (*signUpInfo).Password,
	})

	if err != nil {
		return err
	}
	return nil
}
