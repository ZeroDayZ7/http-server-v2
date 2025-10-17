package router

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/zerodayz7/http-server/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	sessionStore *session.Store,
) {
	SetupHealthRoutes(app)
	setupAuthRoutes(app, authHandler, sessionStore)
	SetupUserRoutes(app, userHandler)
	SetupFallbackHandlers(app)
}
