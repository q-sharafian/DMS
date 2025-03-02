package main

import (
	"DMS/internal/dal"
	"DMS/internal/db"
	repo "DMS/internal/repository"
)

func main() {
	psqlConnDetails := db.PsqlConnDetails{
		Host: "localhost",
		Port: 5192,
		// TODO: Edit username and password such that use enviroment variable
		Username: "mohammad",
		Password: "3522694",
		DB:       "DMS",
	}
	dall := dal.NewPostgresDAL(psqlConnDetails)
	repo.NewRepository(dall)
}
