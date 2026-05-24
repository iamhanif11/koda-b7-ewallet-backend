package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iamhanif11/ewallet-backend/internal/controller"
	"github.com/iamhanif11/ewallet-backend/internal/middleware"
	"github.com/iamhanif11/ewallet-backend/internal/repository"
	"github.com/iamhanif11/ewallet-backend/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TransactionRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRepository := repository.NewAuthRepository(db)
	transactionRepository := repository.NewTransactionRepository()
	authMiddleware := middleware.NewAuthMiddleware(authRepository)
	transactionService := service.NewTransactionService(transactionRepository, db)
	transactionController := controller.NewTransactionController(transactionService)

	TransactionRouter := router.Group("/transaction")

	TransactionRouter.Use(authMiddleware.VerifyToken())

	TransactionRouter.GET("/receivers", transactionController.FindReceivers)

}
