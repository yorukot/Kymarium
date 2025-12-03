package monitor

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/yorukot/knocker/models"
)

// Result captures the outcome of running a monitor.
type Result struct {
	Success  bool
	Duration time.Duration
	Status   models.PingStatus
	Data     any
}

// Run executes a monitor using the default HTTP client.
func Run(ctx context.Context, monitor models.Monitor) (*Result, error) {
	return RunWithClient(ctx, http.DefaultClient, monitor)
}

// RunWithClient executes a monitor with a provided HTTP client (useful for tests).
func RunWithClient(ctx context.Context, client *http.Client, monitor models.Monitor) (*Result, error) {
	if client == nil {
		client = http.DefaultClient
	}

	switch monitor.Type {
	case models.MonitorTypeHTTP:
		return RunHTTP(ctx, client, monitor)
	case models.MonitorTypePing:
		return RunPing(ctx, monitor)
	default:
		return nil, fmt.Errorf("unsupported monitor type %q", monitor.Type)
	}
}
