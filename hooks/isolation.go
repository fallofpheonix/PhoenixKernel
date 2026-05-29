/*
 * REPOSITORY: PhoenixKernel
 * ARCHITECTURAL JUSTIFICATION: Low-level enforcement mechanisms in ring-0.
 * DEPENDENCY BOUNDARY: Direct BPF map manipulation. No AI.
 * DETERMINISTIC CONSIDERATIONS: Absolute authority over process execution.
 */

package hooks

import (
	"fmt"
	"github.com/cilium/ebpf"
)

// Enforcer provides the physical mechanism for isolation.
type Enforcer struct {
	blockedPids *ebpf.Map
}

func NewEnforcer(m *ebpf.Map) *Enforcer {
	return &Enforcer{blockedPids: m}
}

// IsolatePID adds a PID to the kernel blocklist.
func (e *Enforcer) IsolatePID(pid uint32) error {
	var action uint32 = 1
	if err := e.blockedPids.Update(&pid, &action, ebpf.UpdateAny); err != nil {
		return fmt.Errorf("failed to isolate PID %d: %w", pid, err)
	}
	return nil
}

// ReleasePID removes a PID from the kernel blocklist.
func (e *Enforcer) ReleasePID(pid uint32) error {
	if err := e.blockedPids.Delete(&pid); err != nil {
		return fmt.Errorf("failed to release PID %d: %w", pid, err)
	}
	return nil
}

// IsBlocked checks if a PID is currently restricted by the kernel.
func (e *Enforcer) IsBlocked(pid uint32) (bool, error) {
	var action uint32
	if err := e.blockedPids.Lookup(&pid, &action); err != nil {
		if err == ebpf.ErrKeyNotExist {
			return false, nil
		}
		return false, err
	}
	return action == 1, nil
}
