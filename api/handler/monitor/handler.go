package monitor

import "github.com/yorukot/knocker/repository"

// Handler handles monitor-related requests.
type Handler struct {
	Repo repository.Repository
}
