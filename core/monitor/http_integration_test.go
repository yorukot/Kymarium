package monitor

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/models/monitorm"
)

// Integration tests with local mock HTTP server
func TestRunHTTP_Integration(t *testing.T) {
	ctx := context.Background()

	// Create a mock server that handles various test scenarios
	mux := http.NewServeMux()

	// Handler for status code tests
	mux.HandleFunc("/status/200", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/status/503", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service Unavailable"))
	})

	// Handler for delay (used with timeout test)
	mux.HandleFunc("/delay/5", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Handler for redirects
	mux.HandleFunc("/redirect/3", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/redirect/2", http.StatusFound)
	})

	mux.HandleFunc("/redirect/2", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/redirect/1", http.StatusFound)
	})

	mux.HandleFunc("/redirect/1", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/status/200", http.StatusFound)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	tests := []struct {
		name        string
		cfg         monitorm.HTTPMonitorConfig
		wantSuccess bool
		wantStatus  models.PingStatus
	}{
		{
			name: "status 200",
			cfg: monitorm.HTTPMonitorConfig{
				URL:    server.URL + "/status/200",
				Method: monitorm.MethodGet,
			},
			wantSuccess: true,
			wantStatus:  models.PingStatusSuccessful,
		},
		{
			name: "status 503",
			cfg: monitorm.HTTPMonitorConfig{
				URL:    server.URL + "/status/503",
				Method: monitorm.MethodGet,
			},
			wantSuccess: false,
			wantStatus:  models.PingStatusFailed,
		},
		{
			name: "delay timeout",
			cfg: monitorm.HTTPMonitorConfig{
				URL:            server.URL + "/delay/5",
				Method:         monitorm.MethodGet,
				RequestTimeout: 1,
			},
			wantSuccess: false,
			wantStatus:  models.PingStatusTimeout,
		},
		{
			name: "redirect limited",
			cfg: monitorm.HTTPMonitorConfig{
				URL:       server.URL + "/redirect/3",
				Method:    monitorm.MethodGet,
				MaxRedirs: 1,
			},
			wantSuccess: false,
			wantStatus:  models.PingStatusFailed,
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
			// Timeout errors are expected and returned alongside a result
			if err != nil && res == nil {
				t.Fatalf("RunHTTP returned error with no result: %v", err)
			}

			if res == nil {
				t.Fatalf("RunHTTP returned nil result")
			}

			t.Logf("Response: Status=%s, Success=%v, Message=%q, Error=%v", res.Status, res.Success, res.Message, err)

			if res.Status != tt.wantStatus {
				t.Fatalf("expected status %s, got %s (success=%v)", tt.wantStatus, res.Status, res.Success)
			}

			if res.Success != tt.wantSuccess {
				t.Fatalf("expected success=%v, got %v", tt.wantSuccess, res.Success)
			}
		})
	}
}

// TestRunHTTP_HTTPS tests HTTPS with self-signed certificates
func TestRunHTTP_HTTPS(t *testing.T) {
	ctx := context.Background()

	// Create HTTPS server with self-signed cert
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := httptest.NewTLSServer(mux)
	defer server.Close()

	// Client that ignores TLS errors
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	cfg := monitorm.HTTPMonitorConfig{
		URL:            server.URL,
		Method:         monitorm.MethodGet,
		IgnoreTLSError: true,
	}

	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}

	monitor := models.Monitor{
		Type:   models.MonitorTypeHTTP,
		Config: cfgBytes,
	}

	res, err := RunHTTP(ctx, client, monitor)
	if err != nil && res == nil {
		t.Fatalf("RunHTTP returned error with no result: %v", err)
	}

	if res == nil {
		t.Fatalf("RunHTTP returned nil result")
	}

	if res.Status != models.PingStatusSuccessful {
		t.Fatalf("expected status successful, got %s", res.Status)
	}

	if !res.Success {
		t.Fatalf("expected success=true, got %v", res.Success)
	}
}
