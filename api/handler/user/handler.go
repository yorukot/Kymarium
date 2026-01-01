package user

import "github.com/yorukot/knocker/repository"

// Handler handles user-related requests.
type Handler struct {
	Repo repository.Repository
}
