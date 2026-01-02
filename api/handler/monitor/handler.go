package monitor

import "github.com/yorukot/kymarium/repository"

// Handler handles monitor-related requests.
type Handler struct {
	Repo repository.Repository
}
