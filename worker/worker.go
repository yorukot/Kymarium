package worker

import (
	"fmt"
	"strconv"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yorukot/knocker/repository"
	"github.com/yorukot/knocker/utils/config"
	"github.com/yorukot/knocker/worker/handler"
	"github.com/yorukot/knocker/worker/tasks"
	"go.uber.org/zap"
)

// Run starts the background worker processing queues.
func Run(db *pgxpool.Pool) {
	zap.L().Info("Starting worker")
	cfg := config.Env()

	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	redisOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: cfg.RedisPassword,
	}
	regionIDString := strconv.FormatInt(config.RegionByName(cfg.AppRegion).ID, 10)

	queues := map[string]int{
		// Consume only the regional queue for monitor pings plus default for shared tasks (e.g., notifications).
		"critical":     1,
		regionIDString: 6,
		"default":      3,
	}

	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10000,
			Queues:      queues,
		},
	)

	notifier := asynq.NewClient(redisOpt)
	defer notifier.Close()

	repo := repository.New(db)
	h := handler.NewHandler(repo, notifier)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeMonitorPingPattern, h.HandleStartServiceTask)
	mux.HandleFunc(tasks.TypeNotificationDispatch, h.HandleNotificationDispatch)

	if err := srv.Run(mux); err != nil {
		panic(err)
	}
}
