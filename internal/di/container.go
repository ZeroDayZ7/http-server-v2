// internal/di/container.go
package di

import (
	"github.com/zerodayz7/http-server/config"
	"github.com/zerodayz7/http-server/internal/handler"
	"github.com/zerodayz7/http-server/internal/repository/mysql"
	"github.com/zerodayz7/http-server/internal/service"
)

type Container struct {
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler
}

func NewContainer(conn *config.DBConn) *Container {
	userRepo := mysql.NewUserRepository(conn)

	authSvc := service.NewAuthService(userRepo)
	userSvc := service.NewUserService(userRepo)

	return &Container{
		AuthHandler: handler.NewAuthHandler(authSvc),
		UserHandler: handler.NewUserHandler(userSvc),
	}
}
