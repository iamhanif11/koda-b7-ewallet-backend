package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iamhanif11/ewallet-backend/internal/controller"
	"github.com/iamhanif11/ewallet-backend/internal/middleware"
	"github.com/iamhanif11/ewallet-backend/internal/repository"
	"github.com/iamhanif11/ewallet-backend/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func AuthRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := router.Group("/auth")

	authRepository := repository.NewAuthRepository(db, rdb)
	authService := service.NewAuthService(authRepository, db)
	authController := controller.NewAuthController(authService)

	authMiddleware := middleware.NewAuthMiddleware(authRepository)

	authRouter.POST("", authController.Login)
	authRouter.POST("/register", authController.Register)
	authRouter.DELETE("/logout", authMiddleware.VerifyToken(), authController.Logout)
	authRouter.POST("/forgot-password/verify-email", authController.VerifyEmail)
	authRouter.POST("/forgot-password/reset", authController.ResetPassword)
}
