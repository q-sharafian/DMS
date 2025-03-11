package main

import (
	"DMS/internal/controllers"
	"DMS/internal/dal"
	"DMS/internal/db"
	"DMS/internal/logger"
	"DMS/internal/routes"
	"DMS/internal/services"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	lgr := logger.NewSLogger(logger.Debug, nil, os.Stderr)

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
	router.Run(":8080")

}
