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
