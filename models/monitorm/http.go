package monitorm

import (
	"net/http"
)

// HTTPMethod represents the HTTP method used by an HTTP monitor.
type HTTPMethod string

// HTTPMethod values.
const (
	MethodGet     HTTPMethod = http.MethodGet
	MethodPost    HTTPMethod = http.MethodPost
	MethodPut     HTTPMethod = http.MethodPut
	MethodDelete  HTTPMethod = http.MethodDelete
	MethodPatch   HTTPMethod = http.MethodPatch
	MethodHead    HTTPMethod = http.MethodHead
	MethodOptions HTTPMethod = http.MethodOptions
)

// BodyEncoding represents the encoding for HTTP request bodies.
type BodyEncoding string

// BodyEncoding values.
const (
	BodyEncodingJSON BodyEncoding = "json"
	BodyEncodingXML  BodyEncoding = "xml"
)

// HTTPHeader represents a single request header for an HTTP monitor.
type HTTPHeader struct {
	Key   string `json:"key" validate:"required,max=255"`
	Value string `json:"value" validate:"required,max=65535"`
}

// HTTPMonitorConfig represents the expected config shape for HTTP monitors.
// Fields are ordered by importance and functional grouping.
type HTTPMonitorConfig struct {
	// Core request configuration
	URL       string     `json:"url" validate:"required,url"`
	Method    HTTPMethod `json:"method" validate:"oneof=GET POST PUT DELETE PATCH HEAD OPTIONS"`
	MaxRedirs int        `json:"max_redirects" validate:"gte=0,lte=1000"`

	// Request options
	RequestTimeout int               `json:"request_timeout" validate:"gte=0"`
	Headers        []HTTPHeader      `json:"headers,omitempty" validate:"omitempty,dive"`
	BodyEncoding   BodyEncoding      `json:"body_encoding,omitempty" validate:"omitempty,oneof=json xml"`
	Body           string            `json:"body,omitempty" validate:"lte=1000000,omitempty"`

	// Response validation
	UpSideDownMode                bool  `json:"upside_down_mode" validate:"boolean"`
	CertificateExpiryNotification bool  `json:"certificate_expiry_notification" validate:"boolean"`
	IgnoreTLSError                bool  `json:"ignore_tls_error" validate:"boolean"`
	AcceptedStatusCodes           []int `json:"accepted_status_codes" validate:"omitempty,dive,min=100,max=599"`
}
