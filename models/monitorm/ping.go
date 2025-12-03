package monitorm

// PingMonitorConfig represents the config required for a ping (TCP) monitor.
// Fields are ordered by importance and functional grouping.
type PingMonitorConfig struct {
	// Target configuration
	Host string `json:"host" validate:"required,hostname|ip"`

	// Dial options
	TimeoutSeconds int `json:"timeout_seconds" validate:"gte=0"`
	PacketSize     int `json:"packet_size" validate:"omitempty,gte=1,lte=65000"`
}
