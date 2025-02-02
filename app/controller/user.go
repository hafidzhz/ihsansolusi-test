package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/dto"
	"github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/validator"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/logger"
)

type UserController interface {
	RegisterUser(c *fiber.Ctx) error
	Deposit(c *fiber.Ctx) error
}

type userController struct {
	repository repository.UserRepository
}

func NewUserController(userRepository repository.UserRepository) UserController {
	return &userController{
		userRepository,
	}
}

// RegisterUser Register a user.
// @Description Register a user.
// @Summary Register a user
// @Tags User
// @Accept json
// @Produce json
// @Param createUserRequest body dto.CreateUserRequest true "Create new user"
// @Success 200 {object} map[string]string "Account Number"
// @Failure 400,401,403,500 {object} ErrorResponse "Error"
// @Router /daftar [post]
func (controller *userController) RegisterUser(c *fiber.Ctx) error {
	userToCreate := &dto.CreateUserRequest{}

	if err := c.BodyParser(userToCreate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	validate := validator.NewValidator()
	if err := validate.Struct(userToCreate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	createdUser, err := controller.repository.CreateUser(userToCreate)
	if err != nil {
		logger.GetLogger().Errorf("[RegisterUser] failed create user, error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"account_number": createdUser.AccountNumber,
	})
}

// Deposit Deposit money to user account.
// @Description Deposit money to user account.
// @Summary Deposit money to user account
// @Tags User
// @Accept json
// @Produce json
// @Param depositRequest body dto.DepositRequest true "Deposit money to user account"
// @Success 200 {object} dto.User status "Ok
// @Failure 400,401,403,500 {object} ErrorResponse "Error"
// @Router /tabung [post]
func (controller *userController) Deposit(c *fiber.Ctx) error {
	request := &dto.DepositRequest{}

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	validate := validator.NewValidator()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	user, err := controller.repository.DepositToUserAccount(request)
	if err != nil {
		logger.GetLogger().Errorf("[Deposit] failed deposit to user, error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user": dto.ToUser(user),
	})
}

// Withdraw Withdraw money to user account.
// @Description Withdraw money to user account.
// @Summary Withdraw money to user account
// @Tags User
// @Accept json
// @Produce json
// @Param withdrawRequest body dto.WithdrawRequest true "Withdraw money to user account"
// @Success 200 {object} dto.User status "Ok
// @Failure 400,401,403,500 {object} ErrorResponse "Error"
// @Router /tabung [post]
func (controller *userController) Withdraw(c *fiber.Ctx) error {
	request := &dto.DepositRequest{}

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	validate := validator.NewValidator()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	user, err := controller.repository.WithdrawFromUserAccount(request)
	if err != nil {
		logger.GetLogger().Errorf("[Withdraw] failed withdraw from user, error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user": dto.ToUser(user),
	})
}
