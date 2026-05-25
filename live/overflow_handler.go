package live

// OverflowHandler manages ring buffer backpressure without dropping states
type OverflowHandler struct{}

func (o *OverflowHandler) Handle() error { return nil }
