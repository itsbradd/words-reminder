package user

import (
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	path        string
	userHandler Handler
}

func NewRouter(userHandler Handler) Router {
	return Router{path: "/users", userHandler: userHandler}
}

func (r Router) Setup(routerGroup fiber.Router) {
	router := routerGroup.Group(r.path)
	router.Post("/signup", r.userHandler.SignUp)
	router.Post("/login", r.userHandler.Login)
}
