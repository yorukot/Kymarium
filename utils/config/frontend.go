package config

import (
	"net/url"
	"strings"
)

// frontendURL parses the configured frontend URL and normalizes the value we use
// for downstream helpers. It tolerates missing schemes by assuming https.
func frontendURL() (*url.URL, string) {
	raw := strings.TrimSpace(Env().FrontendURL)
	if raw == "" {
		return nil, ""
	}

	// Ensure we have a scheme so url.Parse can separate host/path cleanly.
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, raw
	}

	// Discard any path/query so callers get just the origin components.
	parsed.Path = ""
	parsed.RawPath = ""
	parsed.ForceQuery = false
	parsed.RawQuery = ""
	parsed.Fragment = ""

	return parsed, ""
}

// FrontendOrigin returns the scheme://host[:port] form of the configured
// frontend URL. Falls back to the raw value when parsing fails.
func FrontendOrigin() string {
	parsed, raw := frontendURL()
	if parsed != nil && parsed.Host != "" {
		scheme := parsed.Scheme
		if scheme == "" {
			scheme = "https"
		}
		return scheme + "://" + parsed.Host
	}

	return raw
}

// FrontendDomain returns only the host[:port] portion of the configured
// frontend URL, with any scheme and path stripped.
func FrontendDomain() string {
	parsed, raw := frontendURL()
	if parsed != nil && parsed.Host != "" {
		return parsed.Host
	}

	trimmed := strings.TrimPrefix(strings.TrimPrefix(raw, "https://"), "http://")
	if idx := strings.IndexByte(trimmed, '/'); idx != -1 {
		trimmed = trimmed[:idx]
	}

	return trimmed
}
