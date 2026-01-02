package user

import "github.com/yorukot/kymarium/repository"

// Handler handles user-related requests.
type Handler struct {
	Repo repository.Repository
}
