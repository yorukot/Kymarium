package auth

import (
	"github.com/yorukot/kymarium/repository"
	"github.com/yorukot/kymarium/utils/config"
)

// Handler handles auth-related requests.
type Handler struct {
	Repo        repository.Repository
	OAuthConfig *config.OAuthConfig
}
