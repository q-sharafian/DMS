package main

import (
	"DMS/internal/controllers"
	"DMS/internal/dal"
	"DMS/internal/db"
	"DMS/internal/graph"
	"DMS/internal/hierarchy"
	"DMS/internal/logger"
	"DMS/internal/routes"
	"DMS/internal/services"
	"fmt"
	"os"
	"strconv"
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

	dbPort, err := strconv.Atoi(os.Getenv("PSQL_PORT"))
	if err != nil {
		lgr.Panic(err)
	}
	psqlConnDetails := db.PsqlConnDetails{
		Host:            os.Getenv("PSQL_HOST"),
		Port:            dbPort,
		Username:        os.Getenv("PSQL_USER"),
		Password:        os.Getenv("PSQL_PASSWORD"),
		DB:              os.Getenv("PSQL_DB"),
		MaxConnLifetime: time.Hour,
		MaxIdleConns:    5,
		MAxOpenConns:    5,
	}
	psqlDAL := dal.NewPostgresDAL(psqlConnDetails, lgr, true)

	// Init redis
	dbIndex, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		lgr.Panic(err)
	}
	expireTime, err := strconv.Atoi(os.Getenv("REDIS_EXPIRE"))
	if err != nil {
		lgr.Panic(err)
	}
	redisConnDetails := &db.RedisConnDetails{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbIndex,
		Expire:   time.Second * time.Duration(expireTime),
	}
	redisDAL := dal.NewRedisInMemoeyDAL(redisConnDetails, lgr)
	// Init hierarchy tree
	graphStorage := graph.NewInMemoryDBStorage(redisDAL, []byte("e"), lgr)
	dynamicGraph := graph.NewDynamicGraph(graphStorage, lgr)
	hierarchyTree := hierarchy.NewHierarchyTree(dynamicGraph, lgr)

	simpleService := services.NewSService(&psqlDAL, hierarchyTree, lgr)
	httpController := controllers.NewHttpController(simpleService, lgr)

	router := gin.Default()
	routes.SetupRouter(router, httpController)
	router.Run(fmt.Sprintf(":%s", os.Getenv("GIN_PORT")))

}
