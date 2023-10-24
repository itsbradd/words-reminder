package pkg

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrBodyValidation = errors.New("invalid parameters")
)

var (
	ErrParseReqBody = &FailRes{
		StatusCode: fiber.StatusBadRequest,
		ErrorCode:  fiber.StatusBadRequest,
		Message:    "Invalid parameters",
		Errors:     nil,
	}
)

type SuccessRes[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type FailRes struct {
	StatusCode int    `json:"-"`
	ErrorCode  int    `json:"errorCode"`
	Message    string `json:"message"`
	Errors     any    `json:"errors,omitempty"`
}

func (f *FailRes) Error() string {
	return f.Message
}

func NewBodyValidationErr(err validation.Errors, message ...string) *FailRes {
	msg := ErrBodyValidation.Error()
	if message != nil {
		msg = message[0]
	}
	return &FailRes{
		StatusCode: fiber.StatusBadRequest,
		ErrorCode:  fiber.StatusBadRequest,
		Message:    msg,
		Errors:     err,
	}
}
