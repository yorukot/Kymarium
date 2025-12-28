package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"github.com/yorukot/knocker/api"
	"github.com/yorukot/knocker/db"
	"github.com/yorukot/knocker/schedular"
	"github.com/yorukot/knocker/utils/config"
	"github.com/yorukot/knocker/utils/id"
	"github.com/yorukot/knocker/utils/logger"
	"github.com/yorukot/knocker/worker"
	"go.uber.org/zap"
)

// @title Ridash API
// @version 1.0
// @description This is the Ridash API server for user authentication and management
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// initialize logger
	logger.InitLogger()

	// load environment variables
	_, err := config.InitConfig()
	if err != nil {
		zap.L().Fatal("Error initializing config", zap.Error(err))
	}

	// Initialize sonyflake
	err = id.Init()
	if err != nil {
		zap.L().Fatal("Error initializing snoyflake", zap.Error(err))
	}

	runAll := len(os.Args) < 2 || os.Args[1] == "all"

	pgsql, err := db.InitDatabase()
	if err != nil {
		zap.L().Fatal("Error initializing Postgres", zap.Error(err))
	}
	defer pgsql.Close()
	
	_, err = config.InitRegionConfig(pgsql)
	if err != nil {
		zap.L().Fatal("Error initializing region config", zap.Error(err))
	}
	
	err = config.InitSMTP()
	if err != nil {
		zap.L().Fatal("Error initializing SMTP", zap.Error(err))
	}

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
