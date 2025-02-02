package route

import (
	"github.com/hafidzhz/ihsansolusi-test/app/controller"

	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(
	a *fiber.App,
	userController controller.UserController,
) {
	route := a.Group("")

	route.Post("/daftar", userController.RegisterUser)
	route.Post("/tabung", userController.Deposit)
	route.Post("/tarik", userController.Withdraw)
	route.Get("/saldo/:accountNumber", userController.GetUser)
}
