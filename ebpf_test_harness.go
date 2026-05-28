package kernel

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type MockPublisher struct{}

func (m *MockPublisher) Publish(topic string, event TelemetryEvent) {
	log.Printf("[Harness] Event on %s: %s", topic, event.EventType)
}

func RunHarness() {
	log.Println("=== PHOENIX EBPF TEST HARNESS ===")
	pub := &MockPublisher{}
	loader := NewLoader(pub)
	
	// Assuming the .o file is compiled and present
	err := loader.Load("./src/phoenix_exec.o")
	if err != nil {
		log.Printf("[ERROR] eBPF Load failed: %v", err)
		log.Println("Note: This requires root privileges and a compiled eBPF object.")
		return
	}
	defer loader.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	
	log.Println("[eBPF] Harness: Listening for events... (Press Ctrl+C to stop)")
	<-sig
}
