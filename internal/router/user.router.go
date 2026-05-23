package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iamhanif11/ewallet-backend/internal/controller"
	"github.com/iamhanif11/ewallet-backend/internal/middleware"
	"github.com/iamhanif11/ewallet-backend/internal/repository"
	"github.com/iamhanif11/ewallet-backend/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UserRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRepository := repository.NewAuthRepository(db)
	userRepository := repository.NewUserRepository(db)

	authMiddleware := middleware.NewAuthMiddleware(authRepository)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	UserRouter := router.Group("/user")

	UserRouter.Use(authMiddleware.VerifyToken())

	UserRouter.GET("/profile", userController.GetProfile)
	UserRouter.POST("/profile/pin/check", userController.CheckPin)
	UserRouter.PATCH("/profile", userController.UpdateProfile)
	UserRouter.PATCH("/password", userController.UpdatePassword)
	UserRouter.PATCH("/pin", userController.UpdatePin)
	UserRouter.GET("/wallet", userController.GetDashboardInformation)
	UserRouter.GET("/reports", userController.GetTransactionReport)
}
