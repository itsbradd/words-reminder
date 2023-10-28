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
	ErrParseReqBody = &FailResponse[any]{
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

type FailResponse[T any] struct {
	StatusCode int    `json:"-"`
	ErrorCode  int    `json:"errorCode"`
	Message    string `json:"message"`
	Errors     T      `json:"errors,omitempty"`
}

func (f *FailResponse[T]) GetStatusCode() int {
	return f.StatusCode
}

func (f *FailResponse[T]) GetErrorCode() int {
	return f.ErrorCode
}

func (f *FailResponse[T]) GetErrors() T {
	return f.Errors
}

func (f *FailResponse[T]) Error() string {
	return f.Message
}

func NewBodyValidationErr(err validation.Errors, message ...string) *FailResponse[validation.Errors] {
	msg := ErrBodyValidation.Error()
	if message != nil {
		msg = message[0]
	}
	return &FailResponse[validation.Errors]{
		StatusCode: fiber.StatusBadRequest,
		ErrorCode:  fiber.StatusBadRequest,
		Message:    msg,
		Errors:     err,
	}
}
