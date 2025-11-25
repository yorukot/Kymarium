package api

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func Run(pgsql *pgxpool.Pool) {
	zap.L().Info("Starting API server")
}