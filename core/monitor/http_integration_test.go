package monitor

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/models/monitorm"
)

// Integration tests that exercise real external endpoints (httpstat.us, httpbin, badssl).
func TestRunHTTP_Integration(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		cfg         monitorm.HTTPMonitorConfig
		wantSuccess bool
		wantStatus  models.PingStatus
	}{
		{
			name: "httpbin 200",
			cfg: monitorm.HTTPMonitorConfig{
				URL: "https://httpbin.org/status/200",
			},
			wantSuccess: true,
			wantStatus:  models.PingStatusSuccessful,
		},
		{
			name: "httpbin 503",
			cfg: monitorm.HTTPMonitorConfig{
				URL: "https://httpbin.org/status/503",
			},
			wantSuccess: false,
			wantStatus:  models.PingStatusFailed,
		},
		{
			name: "httpbin delay timeout",
			cfg: monitorm.HTTPMonitorConfig{
				URL:            "https://httpbin.org/delay/5",
				RequestTimeout: 1,
			},
			wantSuccess: false,
			wantStatus:  models.PingStatusTimeout,
		},
		{
			name: "httpbin redirect limited",
			cfg: monitorm.HTTPMonitorConfig{
				URL:       "https://httpbin.org/redirect/3",
				MaxRedirs: 1,
			},
			wantSuccess: false,
			wantStatus:  models.PingStatusFailed,
		},
		{
			name: "badssl self-signed allowed",
			cfg: monitorm.HTTPMonitorConfig{
				URL:            "https://self-signed.badssl.com/",
				IgnoreTLSError: true,
			},
			wantSuccess: true,
			wantStatus:  models.PingStatusSuccessful,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cfgBytes, err := json.Marshal(tt.cfg)
			if err != nil {
				t.Fatalf("marshal config: %v", err)
			}

			monitor := models.Monitor{
				Type:   models.MonitorTypeHTTP,
				Config: cfgBytes,
			}

			runCtx := ctx
			if tt.cfg.RequestTimeout > 0 {
				var cancel context.CancelFunc
				runCtx, cancel = context.WithTimeout(ctx, time.Duration(tt.cfg.RequestTimeout+2)*time.Second)
				defer cancel()
			}

			res, err := RunHTTP(runCtx, http.DefaultClient, monitor)
			if err != nil && tt.wantStatus != models.PingStatusTimeout {
				t.Fatalf("RunHTTP returned error: %v", err)
			}

			if res.Status != tt.wantStatus {
				t.Fatalf("expected status %s, got %s (success=%v)", tt.wantStatus, res.Status, res.Success)
			}

			if res.Success != tt.wantSuccess {
				t.Fatalf("expected success=%v, got %v", tt.wantSuccess, res.Success)
			}
		})
	}
}
