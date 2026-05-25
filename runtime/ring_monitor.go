package runtime

import (
	"github.com/fallofpheonix/phoenix-os/phoenix_os/bus"
)

// RingMonitor tracks pressure and drops in the eBPF ring buffer.
type RingMonitor struct {
	bus *bus.Bus
}

func NewRingMonitor(b *bus.Bus) *RingMonitor {
	return &RingMonitor{bus: b}
}

func (rm *RingMonitor) GetPressure(topic string) float64 {
	return rm.bus.QueuePressure(topic)
}

func (rm *RingMonitor) GetDroppedCount() int64 {
	return rm.bus.Dropped
}
