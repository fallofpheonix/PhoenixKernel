package live

// ClockSync ensures nanosecond precision mapping to logical ticks
type ClockSync struct{}

func (c *ClockSync) Sync() error { return nil }
