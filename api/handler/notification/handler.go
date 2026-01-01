package notification

import "github.com/yorukot/knocker/repository"

// Handler handles notification-related requests.
type Handler struct {
	Repo repository.Repository
}
