package sandbox

import (
	"errors"
	"fmt"
	"time"
)

// KernelSimulator mocks eBPF and kernel runtime constraints.
type KernelEvent struct {
	ID        uint64
	Timestamp uint64
	Data      interface{}
	Size      int
}

type KernelSimulator struct {
	MaxMapEntries      int
	CurrentEntries     int
	StackDepth         int
	MaxStackDepth      int
	GlobalEnergyBudget float64
	ConsumedEnergy     float64
	RingBufferSize     int
	RingBufferUsed     int
	DroppedEvents      uint64
	Events             []KernelEvent
	NextEventID        uint64
}

// NewKernelSimulator initializes a mock kernel environment.
func NewKernelSimulator() *KernelSimulator {
	return &KernelSimulator{
		MaxMapEntries:      1024,
		MaxStackDepth:      512,
		GlobalEnergyBudget: 1000.0,
		ConsumedEnergy:     0.0,
		RingBufferSize:     4096,
		Events:             make([]KernelEvent, 0),
	}
}

// SubmitToRingBuffer simulates writing to an eBPF ring buffer.
func (k *KernelSimulator) SubmitToRingBuffer(data interface{}, size int) error {
	if k.RingBufferUsed+size > k.RingBufferSize {
		k.DroppedEvents++
		return fmt.Errorf("ring buffer overflow: dropped event of size %d", size)
	}
	
	k.NextEventID++
	event := KernelEvent{
		ID:        k.NextEventID,
		Timestamp: uint64(time.Now().UnixNano()),
		Data:      data,
		Size:      size,
	}
	
	k.Events = append(k.Events, event)
	k.RingBufferUsed += size
	return nil
}

// ConsumeFromRingBuffer simulates a userspace agent reading from the ring buffer.
func (k *KernelSimulator) ConsumeFromRingBuffer() (*KernelEvent, error) {
	if len(k.Events) == 0 {
		return nil, errors.New("ring buffer empty")
	}
	
	event := k.Events[0]
	k.Events = k.Events[1:]
	k.RingBufferUsed -= event.Size
	return &event, nil
}

// RequestEnergy attempts to consume energy from the global budget.
func (k *KernelSimulator) RequestEnergy(amount float64) error {
	if k.ConsumedEnergy+amount > k.GlobalEnergyBudget {
		return fmt.Errorf("energy budget exceeded: requested %.2f, available %.2f", amount, k.GlobalEnergyBudget-k.ConsumedEnergy)
	}
	k.ConsumedEnergy += amount
	return nil
}

// UpdateMap simulates writing to an eBPF map.
func (k *KernelSimulator) UpdateMap(key string, value interface{}) error {
	if k.CurrentEntries >= k.MaxMapEntries {
		return errors.New("eBPF map limit reached (memory exhaustion)")
	}
	k.CurrentEntries++
	return nil
}

// CheckStackDepth simulates eBPF verifier stack depth checks.
func (k *KernelSimulator) CheckStackDepth(depth int) error {
	if depth > k.MaxStackDepth {
		return fmt.Errorf("eBPF verifier error: stack depth %d exceeds limit %d", depth, k.MaxStackDepth)
	}
	k.StackDepth = depth
	return nil
}

// Panic simulates a kernel panic.
func (k *KernelSimulator) Panic(reason string) {
	fmt.Printf("!!! KERNEL PANIC: %s !!!\n", reason)
}
