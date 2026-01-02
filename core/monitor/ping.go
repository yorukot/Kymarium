package monitor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	ping "github.com/prometheus-community/pro-bing"
	"github.com/yorukot/kymarium/models"
)

// RunPing executes an ICMP ping monitor using the provided configuration.
func RunPing(ctx context.Context, monitor models.Monitor) (*Result, error) {
	cfg, err := monitor.PingConfig()
	if err != nil {
		return nil, err
	}

	packetSize := cfg.PacketSize
	if packetSize == 0 {
		packetSize = 56
	}

	runPingAttempt := func(privileged bool) (*ping.Pinger, time.Duration, error) {
		pinger, err := ping.NewPinger(cfg.Host)
		if err != nil {
			return nil, 0, fmt.Errorf("create pinger: %w", err)
		}

		pinger.SetPrivileged(privileged)
		pinger.Count = 1
		pinger.Size = packetSize
		pinger.Interval = time.Second
		pinger.Timeout = time.Duration(cfg.TimeoutSeconds) * time.Second
		if pinger.Timeout == 0 {
			pinger.Timeout = 5 * time.Second
		}

		start := time.Now()
		runErr := pinger.RunWithContext(ctx)
		return pinger, time.Since(start), runErr
	}

	pinger, duration, runErr := runPingAttempt(true)
	if isPermissionError(runErr) {
		pinger, duration, runErr = runPingAttempt(false)
	}

	if pinger == nil {
		return nil, runErr
	}

	stats := pinger.Statistics()
	success := stats != nil && stats.PacketsRecv > 0
	status, message := classifyPingOutcome(runErr, success)

	return &Result{
		Success:  success,
		Duration: duration,
		Status:   status,
		Message:  message,
	}, runErr
}

func classifyPingOutcome(runErr error, success bool) (models.PingStatus, string) {
	if success {
		return models.PingStatusSuccessful, ""
	}

	if runErr != nil {
		return classifyPingError(runErr)
	}

	return models.PingStatusTimeout, "no reply received"
}

func classifyPingError(err error) (models.PingStatus, string) {
	if err == nil {
		return models.PingStatusFailed, ""
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return models.PingStatusTimeout, "ping timed out"
	}

	return models.PingStatusFailed, err.Error()
}

func isPermissionError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, os.ErrPermission) {
		return true
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "operation not permitted") || strings.Contains(msg, "permission denied")
}
