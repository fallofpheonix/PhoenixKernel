package kernel

// TelemetryEvent represents a telemetry event emitted by the kernel layer.
type TelemetryEvent struct {
	EventID   string
	SeqID     int64
	Source    string
	EventType string
	Severity  float64
	Payload   []byte
}

// EventPublisher defines how the kernel layer pushes events to the higher-level OS.
type EventPublisher interface {
	Publish(topic string, event TelemetryEvent)
}
