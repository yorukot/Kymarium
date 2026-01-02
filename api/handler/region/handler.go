package region

import "github.com/yorukot/kymarium/repository"

// Handler handles region-related requests.
type Handler struct {
	Repo repository.Repository
}
