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

func TransactionRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRepository := repository.NewAuthRepository(db, rdb)
	transactionRepository := repository.NewTransactionRepository()
	authMiddleware := middleware.NewAuthMiddleware(authRepository)
	transactionService := service.NewTransactionService(transactionRepository, db)
	transactionController := controller.NewTransactionController(transactionService)

	TransactionRouter := router.Group("/transaction")

	TransactionRouter.Use(authMiddleware.VerifyToken())

	TransactionRouter.GET("/receivers", transactionController.FindReceivers)
	TransactionRouter.POST("/transfer", transactionController.Transfer)
	TransactionRouter.POST("/topup", transactionController.Topup)
	TransactionRouter.GET("/history", transactionController.TransactionHistory)

}
