package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iamhanif11/ewallet-backend/internal/controller"
	"github.com/iamhanif11/ewallet-backend/internal/repository"
	"github.com/iamhanif11/ewallet-backend/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRouter := router.Group("/auth")

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository)
	authController := controller.NewAuthController(authService)

	authRouter.POST("", authController.Login)
	authRouter.POST("/register", authController.Register)
}
