package handler

import (
	"github.com/hibiken/asynq"
	"github.com/yorukot/knocker/repository"
)

// Handler coordinates worker task handlers.
type Handler struct {
	repo       repository.Repository
	notifier   *asynq.Client
	pingBuffer *PingRecorder
}

// NewHandler constructs a worker handler with dependencies.
func NewHandler(repo repository.Repository, notifier *asynq.Client) *Handler {
	return &Handler{
		repo:       repo,
		notifier:   notifier,
		pingBuffer: NewPingRecorder(repo),
	}
}
