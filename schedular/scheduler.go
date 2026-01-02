package schedular

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/repository"
	"github.com/yorukot/kymarium/utils/config"
	"github.com/yorukot/kymarium/worker/tasks"
	"go.uber.org/zap"
)

// Run starts the scheduler loop that enqueues monitor pings.
func Run(pgsql *pgxpool.Pool) {
	redisAddr := fmt.Sprintf("%s:%s", config.Env().RedisHost, config.Env().RedisPort)
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: config.Env().RedisPassword,
	})
	defer asynqClient.Close()

	repo := repository.New(pgsql)
	zap.L().Info("Starting scheduler")

	// TODO: Implementing graceful shutdown
	// Create ticker to run every 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		loop(repo, asynqClient)
	}
}

// loop handles a single iteration of fetching and scheduling monitors
func loop(repo repository.Repository, asynqClient *asynq.Client) {
	ctx := context.Background()

	// Start a transaction to fetch monitors
	tx, err := repo.StartTransaction(ctx)
	if err != nil {
		zap.L().Error("Failed to start transaction for fetching monitors", zap.Error(err))
		return
	}
	defer repo.DeferRollback(ctx, tx)

	// Fetch all monitors that need to be pinged
	monitors, err := repo.ListMonitorsDueForCheck(ctx, tx)
	if err != nil {
		zap.L().Error("Failed to fetch monitors due for check", zap.Error(err))
		return
	}

	// Commit the read transaction
	if err := repo.CommitTransaction(ctx, tx); err != nil {
		zap.L().Error("Failed to commit transaction", zap.Error(err))
		return
	}

	if len(monitors) == 0 {
		zap.L().Debug("No monitors due for checking")
		return
	}

	zap.L().Info("Fetched monitors", zap.Int("count", len(monitors)))

	// In this we need to separate the monitors to different goroutines (100-200 monitors per goroutine)
	// then call the scheduleMonitors function to insert into asynq queue
	batchSize := 20 // 100-200 monitors per goroutine
	for i := 0; i < len(monitors); i += batchSize {
		end := i + batchSize
		end = min(end, len(monitors))
		batch := monitors[i:end]

		// Launch goroutine for each batch
		go scheduleMonitors(batch, asynqClient)

		// Batch update last pinged time
		go batchUpdateLastChecked(repo, batch)
	}
}

// batchUpdateLastChecked updates the last_checked and next_check times for a batch of monitors
func batchUpdateLastChecked(repo repository.Repository, monitors []models.Monitor) {
	ctx := context.Background()

	// Start a transaction for updating
	tx, err := repo.StartTransaction(ctx)
	if err != nil {
		zap.L().Error("Failed to start transaction for updating monitors", zap.Error(err))
		return
	}
	defer repo.DeferRollback(ctx, tx)

	now := time.Now()

	// Prepare monitor IDs and their respective next_check times
	monitorIDs := make([]int64, len(monitors))
	nextChecks := make([]time.Time, len(monitors))

	for i, monitor := range monitors {
		monitorIDs[i] = monitor.ID
		// Calculate next check time based on this monitor's specific interval with jitter
		nextChecks[i] = now.Add(time.Duration(monitor.Interval)*time.Second + calculateJitter(monitor.Interval))
	}

	if err := repo.BatchUpdateMonitorsLastChecked(ctx, tx, monitorIDs, nextChecks, now); err != nil {
		zap.L().Error("Failed to batch update monitors last checked time",
			zap.Int("count", len(monitorIDs)),
			zap.Error(err))
		return
	}

	if err := repo.CommitTransaction(ctx, tx); err != nil {
		zap.L().Error("Failed to commit update transaction", zap.Error(err))
		return
	}

	zap.L().Debug("Successfully updated last checked time for monitors",
		zap.Int("count", len(monitors)))
}

// Insert into schedular logic here
// Detail: This basically going insert the monitor task into asynq queue
// Creates one task per monitor per region
func scheduleMonitors(monitors []models.Monitor, asynqClient *asynq.Client) {

	for _, monitor := range monitors {
		// Create a task for each region

		for _, regionID := range monitor.RegionIDs {
			// Create asynq task with region
			task, err := tasks.NewMonitorPing(monitor, regionID)
			if err != nil {
				zap.L().Error("Failed to create monitor task payload",
					zap.Int64("monitor_id", monitor.ID),
					zap.Int64("region_id", regionID),
					zap.Error(err))
				continue
			}

			regionIDString := fmt.Sprintf("%d", regionID)

			// Enqueue the task
			info, err := asynqClient.Enqueue(
				task,
				asynq.Timeout(120*time.Second),
				// Route each region's task to its own queue so only the matching regional worker consumes it.
				asynq.Queue(regionIDString),
			)
			if err != nil {
				zap.L().Error("Failed to enqueue monitor task",
					zap.Int64("monitor_id", monitor.ID),
					zap.Int64("region_id", regionID),
					zap.Error(err))
				continue
			}

			zap.L().Debug("Enqueued monitor task",
				zap.Int64("monitor_id", monitor.ID),
				zap.Int64("region_id", regionID),
				zap.String("task_id", info.ID))
		}
	}
}
