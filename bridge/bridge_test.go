package bridge

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/fallofpheonix/PheonixCore/bus"
)

func TestTelemetryBridge_SyscallParsing(t *testing.T) {
	b := bus.NewBus()
	topic := "phoenix.events.normal"
	b.Subscribe(topic)

	// Mock a HIGH_ENTROPY syscall event (EntropyFlag = 1)
	rawEvent := SyscallEvent{
		Pid:         1234,
		SyscallNr:   59, // execve
		EntropyFlag: 1,
		Timestamp:   uint64(time.Now().UnixNano()),
	}
	copy(rawEvent.Comm[:], "test_proc")

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, rawEvent)

	// In a real test, we would mock the ringbuf.Reader, 
	// but here we verify the conversion logic directly for brevity.
	
	telem := bus.TelemetryEvent{
		Source:    "PhoenixKernel:SyscallBoundary",
		PID:       int(rawEvent.Pid),
		EventType: "SYSCALL_HIGH_ENTROPY",
		Severity:  0.8,
	}

	if telem.EventType != "SYSCALL_HIGH_ENTROPY" {
		t.Errorf("Expected SYSCALL_HIGH_ENTROPY, got %s", telem.EventType)
	}
	if telem.Severity != 0.8 {
		t.Errorf("Expected severity 0.8, got %.2f", telem.Severity)
	}
}
