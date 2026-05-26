package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/iamhanif11/ewallet-backend/internal/config"
	"github.com/iamhanif11/ewallet-backend/internal/router"
	"github.com/joho/godotenv"
)

// @title						Backend Ewallet
// @version						1.0
// @description					Backend created by Hanif using Gin
// @license.name				MIT
// @host						localhost:8080
// @BasePath					/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description					Bearer token used for authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env. \ncause: %s", err.Error())
	}
	app := gin.Default()
	db, err := config.ConnectPsql()
	if err != nil {
		log.Fatalf("Database connection error. \ncause: %s", err.Error())
	}
	defer db.Close()
	log.Printf("Database Connected")

	rc, err := config.ConnectRedis()
	if err != nil {
		log.Fatalf("Redis connection error. \ncause: %s", err.Error())
	}
	defer rc.Close()
	log.Println("Redis Connected")

	router.InitRouter(app, db, rc)

	app.Run(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")))
}
