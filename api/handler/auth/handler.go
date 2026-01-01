package auth

import (
	"github.com/yorukot/knocker/repository"
	"github.com/yorukot/knocker/utils/config"
)

// Handler handles auth-related requests.
type Handler struct {
	Repo        repository.Repository
	OAuthConfig *config.OAuthConfig
}
