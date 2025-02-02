package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/controller"
	"github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/config"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/middleware"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/route"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/logger"
)

func Serve() {
	appCfg := config.AppCfg()

	logger.SetUpLogger()
	logr := logger.GetLogger()

	if err := database.ConnectDB(); err != nil {
		logr.Panicf("failed database setup. error: %v", err)
	}

	fiberCfg := config.FiberConfig()
	app := fiber.New(fiberCfg)

	userRepository := repository.NewUserRepositoryImpl(database.DB)
	userController := controller.NewUserController(userRepository)

	middleware.FiberMiddleware(app)

	route.GeneralRoute(app)
	route.SwaggerRoute(app)
	route.PublicRoutes(app, userController)
	route.NotFoundRoute(app)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sigCh
		logr.Infoln("Shutting down server...")
		_ = app.Shutdown()
	}()

	serverAddr := fmt.Sprintf("%s:%d", appCfg.Host, appCfg.Port)
	if err := app.Listen(serverAddr); err != nil {
		logr.Errorf("Oops... server is not running! error: %v", err)
	}

}
