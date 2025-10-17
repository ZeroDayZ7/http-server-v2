package main

import (
	"os"

	"github.com/zerodayz7/http-server/config"
	"github.com/zerodayz7/http-server/internal/handler"
	mysqlrepo "github.com/zerodayz7/http-server/internal/repository/mysql"
	"github.com/zerodayz7/http-server/internal/router"
	"github.com/zerodayz7/http-server/internal/server"
	"github.com/zerodayz7/http-server/internal/service"
	"github.com/zerodayz7/http-server/internal/shared/logger"
	"go.uber.org/zap"
)

func main() {
	_, _ = logger.InitLogger(os.Getenv("ENV"))

	// config
	if err := config.LoadConfigGlobal(); err != nil {
		logger.GetLogger().Fatal("Config load failed", zap.Error(err))
	}

	// DB
	conn, closeDB := config.MustInitDB()
	defer closeDB()

	// Repos, service, handlers
	userRepo := mysqlrepo.NewUserRepository(conn)
	authSvc := service.NewAuthService(userRepo)
	userSvc := service.NewUserService(userRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)

	// session
	sessionStore := config.InitSessionStore(conn)

	// Fiber
	app := config.NewFiberApp(sessionStore)

	// routes
	router.SetupRoutes(app, authHandler, userHandler, sessionStore)

	// graceful shutdown
	server.SetupGracefulShutdown(app, closeDB, config.AppConfig.Shutdown)

	address := "0.0.0.0:" + config.AppConfig.Server.Port
	logger.GetLogger().Info(
		"Server listening",
		zap.String("address", address),
	)
	if err := app.Listen(address); err != nil {
		logger.GetLogger().Fatal("Failed to start server", zap.Error(err))
	}
}
