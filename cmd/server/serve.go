package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hafidzhz/ihsansolusi-test/app/controller"
	"github.com/hafidzhz/ihsansolusi-test/app/repository"
	"github.com/hafidzhz/ihsansolusi-test/pkg/config"
	"github.com/hafidzhz/ihsansolusi-test/pkg/middleware"
	"github.com/hafidzhz/ihsansolusi-test/pkg/route"
	"github.com/hafidzhz/ihsansolusi-test/platform/database"
	"github.com/hafidzhz/ihsansolusi-test/platform/logger"

	"github.com/gofiber/fiber/v2"
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
