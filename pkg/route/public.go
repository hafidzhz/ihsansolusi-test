package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/controller"
)

func PublicRoutes(
	a *fiber.App,
	userController controller.UserController,
) {
	route := a.Group("")

	route.Post("/daftar", userController.RegisterUser)
	route.Post("/tabung", userController.Deposit)
	// route.Post("/tarik", userController.With)
}
