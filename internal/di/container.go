// internal/di/container.go
package di

import (
	authHandler "github.com/zerodayz7/http-server/internal/features/auth/handler"
	userHandler "github.com/zerodayz7/http-server/internal/features/users/handler"

	authService "github.com/zerodayz7/http-server/internal/features/auth/service"
	userService "github.com/zerodayz7/http-server/internal/features/users/service"

	"github.com/zerodayz7/http-server/internal/features/users/repository/mysql"

	"gorm.io/gorm"
)

// Container przechowuje wszystkie zależności serwisów i handlerów
type Container struct {
	AuthHandler *authHandler.AuthHandler
	UserHandler *userHandler.UserHandler
}

// NewContainer tworzy nowy kontener z wszystkimi zależnościami
func NewContainer(db *gorm.DB) *Container {
	// repozytorium użytkowników
	userRepo := mysql.NewUserRepository(db)

	// serwisy
	authSvc := authService.NewAuthService(userRepo)
	userSvc := userService.NewUserService(userRepo)

	// handlery
	return &Container{
		AuthHandler: authHandler.NewAuthHandler(authSvc),
		UserHandler: userHandler.NewUserHandler(userSvc),
	}
}
