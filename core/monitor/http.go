package monitor

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/models/monitorm"
)

// RunHTTP executes an HTTP monitor using the provided client and monitor config.
func RunHTTP(ctx context.Context, baseClient *http.Client, monitor models.Monitor) (*Result, error) {
	cfg, err := monitor.HTTPConfig()
	if err != nil {
		return nil, err
	}

	method := string(cfg.Method)
	if method == "" {
		method = http.MethodGet
	}

	req, err := http.NewRequestWithContext(ctx, method, cfg.URL, strings.NewReader(cfg.Body))
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}

	if cfg.Body != "" && req.Header.Get("Content-Type") == "" {
		switch strings.ToLower(string(cfg.BodyEncoding)) {
		case "json":
			req.Header.Set("Content-Type", "application/json")
		default:
			req.Header.Set("Content-Type", "text/plain")
		}
	}

	client := prepareHTTPClient(baseClient, cfg)

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)
	if err != nil {
		status := models.PingStatusFailed
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			status = models.PingStatusTimeout
		}
		return &Result{
			Success:    false,
			Duration:   duration,
			StatusCode: status,
			Data: map[string]any{
				"error": err.Error(),
			},
		}, fmt.Errorf("perform http request: %w", err)
	}
	defer resp.Body.Close()

	isAccepted := acceptedStatus(cfg.AcceptedStatusCodes, resp.StatusCode)
	success := isAccepted
	if cfg.UpSideDownMode {
		success = !isAccepted
	}

	status := models.PingStatusFailed
	if success {
		status = models.PingStatusSuccessful
	}

	return &Result{
		Success:    success,
		Duration:   duration,
		StatusCode: status,
		Data: map[string]any{
			"http_status_code": resp.StatusCode,
		},
	}, nil
}

func prepareHTTPClient(base *http.Client, cfg *monitorm.HTTPMonitorConfig) *http.Client {
	client := *base

	applyTimeout(&client, cfg)
	applyRedirects(&client, base, cfg)
	client.Transport = applyTransport(base.Transport, cfg)

	return &client
}

func applyTimeout(client *http.Client, cfg *monitorm.HTTPMonitorConfig) {
	if cfg.RequestTimeout > 0 {
		client.Timeout = time.Duration(cfg.RequestTimeout) * time.Second
	}
}

func applyRedirects(client *http.Client, base *http.Client, cfg *monitorm.HTTPMonitorConfig) {
	if cfg.MaxRedirs <= 0 {
		return
	}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= cfg.MaxRedirs {
			return http.ErrUseLastResponse
		}
		if base.CheckRedirect != nil {
			return base.CheckRedirect(req, via)
		}
		return nil
	}
}

func applyTransport(base http.RoundTripper, cfg *monitorm.HTTPMonitorConfig) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}

	t, ok := base.(*http.Transport)
	if !ok {
		return base
	}

	cloned := t.Clone()
	if cfg.IgnoreTLSError {
		cloned.TLSClientConfig = cloneTLSConfig(cloned.TLSClientConfig)
	}
	return cloned
}

func cloneTLSConfig(conf *tls.Config) *tls.Config {
	if conf == nil {
		return &tls.Config{InsecureSkipVerify: true}
	}

	cloned := conf.Clone()
	cloned.InsecureSkipVerify = true
	return cloned
}

func acceptedStatus(accepted []int, status int) bool {
	if len(accepted) == 0 {
		return status >= 200 && status < 300
	}

	return slices.Contains(accepted, status)
}
