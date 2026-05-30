package runtime

import (
	"github.com/fallofpheonix/PheonixCore/bus"
)

// ReplayBridge connects the replay engine to the live kernel bus.
type ReplayBridge struct {
	liveBus   *bus.Bus
	replayBus *bus.Bus
}

func NewReplayBridge(live, replay *bus.Bus) *ReplayBridge {
	return &ReplayBridge{
		liveBus:   live,
		replayBus: replay,
	}
}

func (rb *ReplayBridge) Sync(event bus.TelemetryEvent) {
	// In a real implementation, this would handle kernel-level replay injection.
	rb.replayBus.Publish("replay_sync", event)
}
