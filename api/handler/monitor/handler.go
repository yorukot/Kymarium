package monitor

import "github.com/jackc/pgx/v5/pgxpool"

type MonitorHandler struct {
	DB *pgxpool.Pool
}
