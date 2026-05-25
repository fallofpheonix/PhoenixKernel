package runtime

import (
	"fmt"
	"github.com/fallofpheonix/phoenix-os/phoenix_os/bus"
)

// ProbeInjector simulates the injection of events into the eBPF ring buffer.
type ProbeInjector struct {
	targetBus *bus.Bus
}

func NewProbeInjector(b *bus.Bus) *ProbeInjector {
	return &ProbeInjector{targetBus: b}
}

func (pi *ProbeInjector) InjectEvent(topic string, event bus.TelemetryEvent) {
	pi.targetBus.Publish(topic, event)
}

func (pi *ProbeInjector) StressBurst(topic string, count int, severity float64) {
	for i := 0; i < count; i++ {
		pi.targetBus.Publish(topic, bus.TelemetryEvent{
			SeqID:    int64(i),
			Severity: severity,
		})
	}
}
