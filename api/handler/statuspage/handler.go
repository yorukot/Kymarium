package statuspage

import "github.com/yorukot/knocker/repository"

// Handler contains dependencies for status page endpoints.
type Handler struct {
	Repo repository.Repository
}
