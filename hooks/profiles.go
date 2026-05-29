/*
 * REPOSITORY: PhoenixKernel
 * ARCHITECTURAL JUSTIFICATION: Management of per-process syscall profiles.
 * DEPENDENCY BOUNDARY: BPF map manipulation.
 * DETERMINISTIC CONSIDERATIONS: Atomic updates to bitmask structures.
 */

package hooks

import (
	"fmt"
	"github.com/cilium/ebpf"
)

// SyscallProfile matches the C struct profile in syscalls.c
type SyscallProfile struct {
	Bitmask [8]uint64 // 512 bits / 64
}

// ProfileManager handles the lifecycle of syscall allow-lists.
type ProfileManager struct {
	profilesMap *ebpf.Map
}

func NewProfileManager(m *ebpf.Map) *ProfileManager {
	return &ProfileManager{profilesMap: m}
}

// LoadProfile pushes a syscall bitmask for a specific PID into the kernel.
func (p *ProfileManager) LoadProfile(pid uint32, profile SyscallProfile) error {
	if err := p.profilesMap.Update(&pid, &profile, ebpf.UpdateAny); err != nil {
		return fmt.Errorf("failed to load syscall profile for PID %d: %w", pid, err)
	}
	return nil
}

// ClearProfile removes a syscall profile from the kernel.
func (p *ProfileManager) ClearProfile(pid uint32) error {
	if err := p.profilesMap.Delete(&pid); err != nil {
		return fmt.Errorf("failed to clear syscall profile for PID %d: %w", pid, err)
	}
	return nil
}

// SetSyscallBit updates a single bit in the profile for a given syscall number.
func (p *ProfileManager) SetSyscallBit(profile *SyscallProfile, syscallNr uint32, allowed bool) {
	if syscallNr >= 512 {
		return
	}
	wordIdx := syscallNr / 64
	bitIdx := syscallNr % 64
	if allowed {
		profile.Bitmask[wordIdx] |= (1 << bitIdx)
	} else {
		profile.Bitmask[wordIdx] &= ^(1 << bitIdx)
	}
}
