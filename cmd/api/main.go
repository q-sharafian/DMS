package main

import (
	"DMS/internal/controllers"
	"DMS/internal/dal"
	"DMS/internal/db"
	"DMS/internal/logger"
	"DMS/internal/routes"
	"DMS/internal/services"
	"fmt"
	"os"
	"time"

	_ "DMS/docs/api"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @version         1.0
// @description     Documentation for DMS API

// @contact.name   Qasem Sharafian

// @license.name  Commercial License

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	var appMode = os.Getenv("APP_MODE")
	if appMode != "production" {
		if err := godotenv.Load(".env"); err != nil {
			panic(fmt.Sprintf("Error loading .env file: %s", err.Error()))
		}
	}
	lgr := logger.NewSLogger(logger.Debug, nil, os.Stderr)
	if appMode != "production" {
		if jwtPrivate, err := os.ReadFile(os.Getenv("JWT_PRIVATE_KEY_FILE_PATH")); err != nil {
			lgr.Panic(err)
		} else if err = os.Setenv("JWT_PRIVATE_KEY", string(jwtPrivate)); err != nil {
			lgr.Panic(err)
		}
		if jwtPublic, err := os.ReadFile(os.Getenv("JWT_PUBLIC_KEY_FILE_PATH")); err != nil {
			lgr.Panic(err)
		} else if err = os.Setenv("JWT_PUBLIC_KEY", string(jwtPublic)); err != nil {
			lgr.Panic(err)
		}
	}

	psqlConnDetails := db.PsqlConnDetails{
		Host: "localhost",
		Port: 5432,
		// TODO: Edit username and password such that use enviroment variable
		Username:        "mohammad",
		Password:        "3522694",
		DB:              "dms",
		MaxConnLifetime: time.Hour,
		MaxIdleConns:    5,
		MAxOpenConns:    5,
	}
	psqlDAL := dal.NewPostgresDAL(psqlConnDetails, lgr, true)
	simpleService := services.NewSService(&psqlDAL, lgr)
	httpController := controllers.NewHttpController(simpleService, lgr)

	router := gin.Default()
	routes.SetupRouter(router, httpController)
	router.Run(fmt.Sprintf(":%s", os.Getenv("GIN_PORT")))

}
