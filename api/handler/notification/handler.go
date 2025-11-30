package notification

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationHandler struct {
	DB *pgxpool.Pool
}
