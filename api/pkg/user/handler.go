package user

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/sonngocme/words-reminder-be/db"
)

type Service interface {
	SignUpUser(context.Context, db.SignUpUserParams) (int64, error)
	HashPassword(string) (string, error)
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

	_, err = h.s.SignUpUser(context.Background(), db.SignUpUserParams{
		Username: (*signUpInfo).Username,
		Password: hashedPass,
	})

	if err != nil {
		return err
	}
	return nil
}
