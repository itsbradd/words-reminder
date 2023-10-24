package pkg

import "github.com/gofiber/fiber/v2"

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
