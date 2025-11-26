package schedular

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/utils/config"
	"github.com/yorukot/knocker/worker/tasks"
	"go.uber.org/zap"
)

func Run(pgsql *pgxpool.Pool) {
	redisAddr := fmt.Sprintf("%s:%s", config.Env().RedisHost, config.Env().RedisPort)
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: config.Env().RedisPassword,
	})
	defer asynqClient.Close()
	zap.L().Info("Starting scheduler")

	// TODO: Implementing graceful shutdown
	// Create ticker to run every 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		loop(pgsql, asynqClient)
	}
}

// loop handles a single iteration of fetching and scheduling monitors
func loop(pgsql *pgxpool.Pool, asynqClient *asynq.Client) {
	// ctx := context.Background()
	// first we need to fetch all monitors that need to be pinged
	monitors := []models.Monitor{}
	// monitors, err := repository.FetchMonitor(ctx, pgsql)
	// if err != nil {
		// zap.L().Error("Failed to fetch monitors", zap.Error(err))
		// return
	// }
	zap.L().Info("Fetched monitors", zap.Int("count", len(monitors)))

	// In this we need to saparate the monitors to the different goroutines it should be 100-200 monitor per goroutine
	// then call the scheduleMonitors function to insert into asynq queue
	batchSize := 150 // 100-200 monitors per goroutine
	for i := 0; i < len(monitors); i += batchSize {
		end := i + batchSize
		end = min(end, len(monitors))
		batch := monitors[i:end]

		// Launch goroutine for each batch
		go scheduleMonitors(batch, asynqClient)
	}
}

// Insert into schedular logic here
// Detail: This basically going insert the monitor task into asynq queue
func scheduleMonitors(monitors []models.Monitor, asynqClient *asynq.Client) {
	for _, monitor := range monitors {
		// Create asynq task
		task, err := tasks.NewMonitorPing(monitor)

		// Enqueue the task
		info, err := asynqClient.Enqueue(task)
		if err != nil {
			zap.L().Error("Failed to enqueue monitor task",
				zap.Int64("monitor_id", monitor.ID),
				zap.String("url", monitor.URL),
				zap.Error(err))
			continue
		}

		zap.L().Debug("Enqueued monitor task",
			zap.Int64("monitor_id", monitor.ID),
			zap.String("url", monitor.URL),
			zap.String("task_id", info.ID))
	}
}
