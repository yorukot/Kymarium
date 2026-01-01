package region

import "github.com/yorukot/knocker/repository"

// Handler handles region-related requests.
type Handler struct {
	Repo repository.Repository
}
