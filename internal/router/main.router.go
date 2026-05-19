package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamhanif11/ewallet-backend/internal/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool) {
	//middleware global
	//router.Use(middleware.Logger)
	//router.Use(middleware.CORSMiddleware)

	AuthRouter(router, db)
	UserRouter(router, db)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dto.Response{
			Message: "Invalid Route",
			Success: false,
			Error:   "Not Found",
		})
	})
}
