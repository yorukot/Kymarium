package team

import "github.com/yorukot/knocker/repository"

// Handler handles team-related requests.
type Handler struct {
	Repo repository.Repository
}
