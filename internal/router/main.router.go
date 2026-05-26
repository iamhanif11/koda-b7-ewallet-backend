package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/iamhanif11/ewallet-backend/docs"
	"github.com/iamhanif11/ewallet-backend/internal/dto"
	"github.com/iamhanif11/ewallet-backend/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	//middleware global
	router.Use(middleware.CORSMiddleware)
	//swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Static("/img/profile", "public/img/profiles")

	AuthRouter(router, db, rdb)
	UserRouter(router, db, rdb)
	TransactionRouter(router, db, rdb)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "Invalid Route",
			Success: false,
			Error:   "Not Found",
		})
	})
}
