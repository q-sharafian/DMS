package main

import (
	"DMS/internal/controllers"
	"DMS/internal/dal"
	"DMS/internal/db"
	"DMS/internal/logger"
	"DMS/internal/routes"
	"DMS/internal/services"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	lgr := logger.NewSLogger(logger.Debug, nil, os.Stderr)

	psqlConnDetails := db.PsqlConnDetails{
		Host: "localhost",
		Port: 5192,
		// TODO: Edit username and password such that use enviroment variable
		Username: "mohammad",
		Password: "3522694",
		DB:       "DMS",
	}
	psqlDAL := dal.NewPostgresDAL(psqlConnDetails, lgr)
	simpleService := services.NewsService(&psqlDAL, lgr)
	httpController := controllers.NewHttpController(simpleService)

	router := gin.Default()
	routes.SetupRouter(router, httpController)
	router.Run(":8080")

}
