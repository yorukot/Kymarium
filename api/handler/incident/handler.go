package incident

import "github.com/yorukot/knocker/repository"

// Handler groups dependencies for incident endpoints.
// Handler handles incident-related requests.
type Handler struct {
	Repo repository.Repository
}
