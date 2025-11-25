package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"github.com/yorukot/knocker/api"
	"github.com/yorukot/knocker/db"
	"github.com/yorukot/knocker/helpers/config"
	"github.com/yorukot/knocker/helpers/logger"
	"github.com/yorukot/knocker/schedular"
	"github.com/yorukot/knocker/worker"
	"go.uber.org/zap"
)

// main is the entry point of the application. It checks command-line arguments to determine which components to run:
// API server, worker, scheduler, or all of them if no specific argument is provided.
func main() {
	// initialize logger
	logger.InitLogger()

	// load environment variables
	_, err := config.InitConfig()
	if err != nil {
		zap.L().Fatal("Error initializing config", zap.Error(err))
	}

	runAll := len(os.Args) < 2 || os.Args[1] == "all"

	pgsql, err := db.InitDatabase()
	if err != nil {
		zap.L().Fatal("Error initializing Postgres", zap.Error(err))
	}
	defer pgsql.Close()

	if runAll || os.Args[1] == "api" {
		go api.Run(pgsql)
	}

	if runAll || os.Args[1] == "worker" {
		go worker.Run(pgsql)
	}

	if runAll || os.Args[1] == "schedular" {
		go schedular.Run(pgsql)
	}

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("Shutting down gracefully...")
}
