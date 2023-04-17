package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := initLogger()
	database := initDatabase("dsn", logger)

	profileRepository := repositories.NewProfileRepository(database, logger)
	profileService := services.NewProfileService(&profileRepository)
	profileController := controllers.NewProfileController(&profileService, logger)

	restServer := servers.NewRESTServer("address", nil, &profileController, logger) // TODO: need to add address

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		restServer.ListenAndServe()
	}()

	go func() {
		defer wg.Done()
		<-ctx.Done()

		restServer.Shutdown(context.Background())
		logger.Info("REST server closed")

		_ = database.Close()
		logger.Info("database closed")
	}()
	wg.Wait()
}
