package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yorukot/knocker/utils/config"
)

type AuthHandler struct {
	DB          *pgxpool.Pool
	OAuthConfig *config.OAuthConfig
}
