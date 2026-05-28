package kernel

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cilium/ebpf/ringbuf"
)

// Enforcer reads eBPF ring buffer events and ingests them into the Phoenix Matrix Event Bus.
type Enforcer struct {
	Pub    EventPublisher
	Reader *ringbuf.Reader
}

// NewEnforcer initializes the observer layer.
func NewEnforcer(pub EventPublisher, r *ringbuf.Reader) *Enforcer {
	return &Enforcer{
		Pub:    pub,
		Reader: r,
	}
}

// Observe starts the blocking loop that reads from the eBPF map and publishes to the Bus.
func (e *Enforcer) Observe(ctx context.Context) {
	log.Println("[eBPF Enforcer] Observer Layer active. Wiring telemetry to Event Bus...")
	
	if e.Reader == nil {
		log.Println("[eBPF Enforcer] Running in Mock/Dry-Run mode (No eBPF Reader provided)")
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("[eBPF Enforcer] Shutting down Observer Layer.")
			return
		default:
			record, err := e.Reader.Read()
			if err != nil {
				log.Printf("[eBPF Enforcer] Ring buffer read error: %v", err)
				continue
			}

			eventID := generateEventID()
			payload, _ := json.Marshal(map[string]interface{}{
				"source_probe": "sys_enter_execve",
				"data_len":     len(record.RawSample),
			})
			
			// The LamportClock and Hash are formally stamped by the Bus during Publish.
			event := TelemetryEvent{
				EventID:   eventID,
				EventType: "syscall_execve",
				Source:    "ebpf_enforcer",
				Severity:  0.8, // execve is a high-severity observable
				Payload:   payload,
			}

			// Ship to the matrix. This guarantees Causal Monotonicity.
			e.Pub.Publish("telemetry", event)
		}
	}
}

func generateEventID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("evt-%x", b)
}
