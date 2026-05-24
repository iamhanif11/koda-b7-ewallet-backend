package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/iamhanif11/ewallet-backend/docs"
	"github.com/iamhanif11/ewallet-backend/internal/dto"
	"github.com/iamhanif11/ewallet-backend/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool) {
	//middleware global
	router.Use(middleware.CORSMiddleware)
	//swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	AuthRouter(router, db)
	UserRouter(router, db)
	TransactionRouter(router, db)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "Invalid Route",
			Success: false,
			Error:   "Not Found",
		})
	})
}
