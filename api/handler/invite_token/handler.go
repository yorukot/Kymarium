package invite_token

import "github.com/yorukot/knocker/repository"

// InviteTokenHandler handles invite token routes.
type InviteTokenHandler struct {
	Repo repository.Repository
}
