package monitor

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/models/monitorm"
)

// RunHTTP executes an HTTP monitor using the provided client and monitor config.
func RunHTTP(ctx context.Context, baseClient *http.Client, monitor models.Monitor) (*Result, error) {
	cfg, err := monitor.HTTPConfig()
	if err != nil {
		return nil, err
	}

	method := string(cfg.Method)

	req, err := http.NewRequestWithContext(ctx, method, cfg.URL, strings.NewReader(cfg.Body))
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	for _, header := range cfg.Headers {
		req.Header.Set(header.Key, header.Value)
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
		status, message := classifyHTTPError(err)
		return &Result{
			Success:  false,
			Duration: duration,
			Status:   status,
			Message:  message,
		}, fmt.Errorf("%s: %w", message, err)
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

	message := ""
	if !success {
		message = statusErrorMessage(resp.StatusCode, cfg.UpSideDownMode)
	}

	return &Result{
		Success:  success,
		Duration: duration,
		Status:   status,
		Message:  message,
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

func classifyHTTPError(err error) (models.PingStatus, string) {
	if err == nil {
		return models.PingStatusFailed, ""
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return models.PingStatusTimeout, "request timed out"
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return models.PingStatusTimeout, "request timed out"
	}

	if tlsMsg, ok := tlsErrorMessage(err); ok {
		return models.PingStatusFailed, tlsMsg
	}

	return models.PingStatusFailed, err.Error()
}

func tlsErrorMessage(err error) (string, bool) {
	var certInvalid *x509.CertificateInvalidError
	if errors.As(err, &certInvalid) {
		return fmt.Sprintf("tls certificate invalid: %v", certInvalid), true
	}

	var unknownAuth x509.UnknownAuthorityError
	if errors.As(err, &unknownAuth) {
		return fmt.Sprintf("tls unknown authority: %v", unknownAuth), true
	}

	var hostnameErr x509.HostnameError
	if errors.As(err, &hostnameErr) {
		return fmt.Sprintf("tls hostname mismatch: %v", hostnameErr), true
	}

	var systemErr x509.SystemRootsError
	if errors.As(err, &systemErr) {
		return fmt.Sprintf("tls root validation failed: %v", systemErr), true
	}

	var recordErr *tls.RecordHeaderError
	if errors.As(err, &recordErr) {
		return fmt.Sprintf("tls handshake failed: %v", recordErr), true
	}

	return "", false
}

func statusErrorMessage(statusCode int, upsideDown bool) string {
	statusText := http.StatusText(statusCode)
	if statusText != "" {
		statusText = " " + statusText
	}

	if upsideDown {
		return fmt.Sprintf("upside_down_mode: received HTTP %d%s which counts as failure", statusCode, statusText)
	}

	return fmt.Sprintf("received HTTP %d%s", statusCode, statusText)
}
