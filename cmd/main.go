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

	router.InitRouter(app, db)

	app.Run(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")))
}
