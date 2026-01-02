package notification

import "github.com/yorukot/kymarium/repository"

// Handler handles notification-related requests.
type Handler struct {
	Repo repository.Repository
}
