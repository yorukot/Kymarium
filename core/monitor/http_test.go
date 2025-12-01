package monitor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/models/monitorm"
)

func TestRunHTTP_ReportsHTTPStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	monitor := monitorFromConfig(t, monitorm.HTTPMonitorConfig{URL: server.URL})

	res, err := RunHTTP(context.Background(), http.DefaultClient, monitor)
	if err != nil {
		t.Fatalf("RunHTTP returned error: %v", err)
	}

	if res.Success {
		t.Fatalf("expected success=false, got true")
	}

	data := resultData(t, res.Data)
	msg := data["error"].(string)
	if !strings.Contains(msg, "500") {
		t.Fatalf("expected error message to mention status code, got %q", msg)
	}
}

func TestRunHTTP_TimeoutErrorMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
	}))
	defer server.Close()

	client := &http.Client{Timeout: 50 * time.Millisecond}
	monitor := monitorFromConfig(t, monitorm.HTTPMonitorConfig{URL: server.URL})

	res, err := RunHTTP(context.Background(), client, monitor)
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}

	if res.Status != models.PingStatusTimeout {
		t.Fatalf("expected timeout status, got %s", res.Status)
	}

	data := resultData(t, res.Data)
	msg := strings.ToLower(data["error"].(string))
	if !strings.Contains(msg, "timeout") && !strings.Contains(msg, "timed out") {
		t.Fatalf("expected timeout message, got %q", msg)
	}
}

func TestRunHTTP_TLSErrorMessage(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	monitor := monitorFromConfig(t, monitorm.HTTPMonitorConfig{URL: server.URL})

	res, err := RunHTTP(context.Background(), http.DefaultClient, monitor)
	if err == nil {
		t.Fatalf("expected TLS error, got nil")
	}

	if res == nil {
		t.Fatalf("expected result even when error occurs")
	}

	data := resultData(t, res.Data)
	msg := strings.ToLower(data["error"].(string))
	if !strings.Contains(msg, "tls") && !strings.Contains(msg, "certificate") {
		t.Fatalf("expected TLS-related message, got %q", msg)
	}
}

func monitorFromConfig(t *testing.T, cfg monitorm.HTTPMonitorConfig) models.Monitor {
	t.Helper()

	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}

	return models.Monitor{
		Type:   models.MonitorTypeHTTP,
		Config: cfgBytes,
	}
}

func resultData(t *testing.T, data any) map[string]any {
	t.Helper()

	m, ok := data.(map[string]any)
	if !ok {
		t.Fatalf("expected map data, got %T", data)
	}
	return m
}
