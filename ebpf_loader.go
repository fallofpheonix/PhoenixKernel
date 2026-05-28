package kernel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
)

// ExecEvent matches the C struct in phoenix_exec.c
type ExecEvent struct {
	Pid        uint32
	Ppid       uint32
	Tgid       uint32
	NsproxyIno uint32
	Uid        uint32
	Comm       [16]byte
	Filename   [128]byte
}

type Loader struct {
	coll  *ebpf.Collection
	links []link.Link
	Pub   EventPublisher
	rb    *ringbuf.Reader
	blockedPids *ebpf.Map
}

func NewLoader(pub EventPublisher) *Loader {
	return &Loader{Pub: pub}
}

func (l *Loader) Load(path string) error {
	spec, err := ebpf.LoadCollectionSpec(path)
	if err != nil {
		return fmt.Errorf("failed to load ebpf spec: %w", err)
	}

	l.coll, err = ebpf.NewCollection(spec)
	if err != nil {
		return fmt.Errorf("failed to create ebpf collection: %w", err)
	}

	// 1. Attach Tracepoint
	prog := l.coll.Programs["handle_execve"]
	if prog == nil {
		return fmt.Errorf("program handle_execve not found")
	}

	tp, err := link.Tracepoint("syscalls", "sys_enter_execve", prog, nil)
	if err != nil {
		return fmt.Errorf("failed to attach tracepoint: %w", err)
	}
	l.links = append(l.links, tp)

	// 2. Attach LSM Program (Reflexive Actuation)
	lsmProg := l.coll.Programs["phoenix_enforce_exec"]
	if lsmProg != nil {
		lsmLink, err := link.AttachLSM(link.LSMOptions{Program: lsmProg})
		if err != nil {
			log.Printf("[eBPF] Warning: Failed to attach LSM program (Is BPF_LSM enabled?): %v", err)
		} else {
			l.links = append(l.links, lsmLink)
			log.Println("[eBPF] Loader: LSM program phoenix_enforce_exec attached.")
		}
	}

	// 2b. Attach mprotect LSM Program
	mproProg := l.coll.Programs["phoenix_enforce_mprotect"]
	if mproProg != nil {
		mproLink, err := link.AttachLSM(link.LSMOptions{Program: mproProg})
		if err != nil {
			log.Printf("[eBPF] Warning: Failed to attach mprotect LSM program: %v", err)
		} else {
			l.links = append(l.links, mproLink)
			log.Println("[eBPF] Loader: LSM program phoenix_enforce_mprotect attached.")
		}
	}

	// 3. Setup Blocked PIDs Map
	l.blockedPids = l.coll.Maps["blocked_pids"]
	if l.blockedPids == nil {
		return fmt.Errorf("map blocked_pids not found")
	}

	rbMap := l.coll.Maps["rb"]
	l.rb, err = ringbuf.NewReader(rbMap)
	if err != nil {
		return fmt.Errorf("failed to create ringbuf reader: %w", err)
	}

	go l.pollEvents()

	log.Println("[eBPF] Loader: Telemetry and Enforcement layers active.")
	return nil
}

// BlockPID adds a PID to the kernel-space blocklist.
func (l *Loader) BlockPID(pid uint32) error {
	if l.blockedPids == nil {
		log.Printf("[eBPF] Warning: blocked_pids map not initialized. Skipping kernel-level block for PID %d.", pid)
		return nil
	}
	var action uint32 = 1
	return l.blockedPids.Update(pid, action, ebpf.UpdateAny)
}

// UnblockPID removes a PID from the kernel-space blocklist.
func (l *Loader) UnblockPID(pid uint32) error {
	if l.blockedPids == nil {
		return nil
	}
	return l.blockedPids.Delete(pid)
}

func (l *Loader) pollEvents() {
	for {
		rec, err := l.rb.Read()
		if err != nil {
			log.Printf("[eBPF] Poll error: %v", err)
			return
		}

		var event ExecEvent
		if err := binary.Read(bytes.NewReader(rec.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("[eBPF] Failed to decode event: %v", err)
			continue
		}

		telemetryEvent := TelemetryEvent{
			SeqID:     int64(event.Pid), // Mapping PID for now, or add counter
			Source:    "ebpf:exec",
			EventType: "EXEC_EVENT",
			Payload:   event.Filename[:],
		}
		if l.Pub != nil {
			l.Pub.Publish("exec", telemetryEvent)
		}
	}
}

func (l *Loader) Close() {
	if l.rb != nil {
		l.rb.Close()
	}
	for _, link := range l.links {
		link.Close()
	}
	if l.coll != nil {
		l.coll.Close()
	}
	log.Println("[eBPF] Loader: Resources released.")
}
