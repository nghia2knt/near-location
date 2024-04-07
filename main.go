package main

import (
	"context"
	"near-location/internal/controller"
	"near-location/internal/repo"
	"near-location/internal/service"
	"near-location/pkg/config"
	"near-location/pkg/database"
	"near-location/router"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	// init config
	log.Info("start init config")
	err := config.InitConfig("./configs", os.Getenv("STAGE"))
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}
	log.Info(config.CV)

	// init database
	log.Info("start init mongo")
	locationDB, err := database.NewMongoDBConnection(ctx, config.CV.MongoConfig.DatabaseLocation)
	if err != nil {
		log.Fatal(err)
	}

	// init internal
	repo := repo.NewRepo(locationDB)
	userService := service.NewUserService(repo)
	controller := controller.NewController(userService)

	router.NewRouter(controller)
}
