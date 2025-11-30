package worker

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yorukot/knocker/utils/config"
	"github.com/yorukot/knocker/worker/handler"
	"github.com/yorukot/knocker/worker/tasks"
	"go.uber.org/zap"
)

func Run(db *pgxpool.Pool) {
	zap.L().Info("Starting worker")
	cfg := config.Env()

	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: cfg.RedisPassword,
		},
		asynq.Config{
			Concurrency: 10000,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	h := handler.NewHandler(db)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeMonitorPing, h.HandleStartServiceTask)

	if err := srv.Run(mux); err != nil {
		panic(err)
	}
}
