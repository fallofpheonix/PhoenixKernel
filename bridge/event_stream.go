/*
 * REPOSITORY: PhoenixKernel
 * ARCHITECTURAL JUSTIFICATION: Telemetry bridge from eBPF ring buffer to PheonixCore Event Bus.
 * DEPENDENCY BOUNDARY: Depends on PheonixCore/bus. Telemetry only.
 * DETERMINISTIC CONSIDERATIONS: Non-blocking dispatch, monotonic time anchoring.
 */

package bridge

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"time"

	"github.com/cilium/ebpf/ringbuf"
	"github.com/fallofpheonix/PheonixCore/bus"
)

// SyscallEvent must match the C struct in syscalls.c
type SyscallEvent struct {
	Pid         uint32
	SyscallNr   uint32
	EntropyFlag uint32
	Timestamp   uint64
	Comm        [16]byte
}

// TelemetryBridge orchestrates the kernel-to-bus event flow.
type TelemetryBridge struct {
	Bus *bus.Bus
	rb  *ringbuf.Reader
}

func NewTelemetryBridge(b *bus.Bus, r *ringbuf.Reader) *TelemetryBridge {
	return &TelemetryBridge{
		Bus: b,
		rb:  r,
	}
}

// Start consumes events from the ring buffer and dispatches them to the bus.
func (t *TelemetryBridge) Start() {
	log.Println("[PhoenixKernel] Telemetry Bridge started.")
	for {
		record, err := t.rb.Read()
		if err != nil {
			log.Printf("[PhoenixKernel] Ring buffer error: %v", err)
			return
		}

		var event SyscallEvent
		if err := binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("[PhoenixKernel] Failed to decode kernel event: %v", err)
			continue
		}

		// Convert to canonical TelemetryEvent
		payload, _ := json.Marshal(event)
		eventType := "SYSCALL_STABLE"
		severity := 0.1
		if event.EntropyFlag == 1 {
			eventType = "SYSCALL_HIGH_ENTROPY"
			severity = 0.8
		}

		telem := bus.TelemetryEvent{
			SeqID:        int64(event.Timestamp),
			MonotonicNs:  int64(event.Timestamp),
			WallTimeUnix: time.Now().Unix(),
			Source:       "PhoenixKernel:SyscallBoundary",
			PID:          int(event.Pid),
			EventType:    eventType,
			Severity:     severity,
			Payload:      payload,
		}

		// Non-blocking publish to the central bus
		go t.Bus.Publish("phoenix.events.normal", telem)
	}
}
