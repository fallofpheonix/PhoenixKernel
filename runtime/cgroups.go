/*
 * REPOSITORY: PheonixKernel
 * ARCHITECTURAL JUSTIFICATION: Cgroup V2 manipulation for process freezing.
 * DEPENDENCY BOUNDARY: Linux kernel interfaces. No dependencies on strategic layers.
 * DETERMINISTIC CONSIDERATIONS: File-based state control for process execution.
 */

package runtime

import (
	"fmt"
	"os"
	"path/filepath"
)

const CGROUP_BASE = "/sys/fs/cgroup/phoenix_isolation"

// FreezePID moves a process into the isolation cgroup and halts its execution.
func FreezePID(pid int) error {
	groupPath := filepath.Join(CGROUP_BASE, fmt.Sprintf("proc_%d", pid))
	
	// 1. Create the cgroup for this process
	if err := os.MkdirAll(groupPath, 0755); err != nil {
		return fmt.Errorf("failed to create isolation cgroup: %w", err)
	}

	// 2. Add the PID to the cgroup
	procsFile := filepath.Join(groupPath, "cgroup.procs")
	if err := os.WriteFile(procsFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		return fmt.Errorf("failed to move PID %d to cgroup: %w", err)
	}

	// 3. Trigger the freeze
	freezeFile := filepath.Join(groupPath, "cgroup.freeze")
	if err := os.WriteFile(freezeFile, []byte("1"), 0644); err != nil {
		return fmt.Errorf("failed to freeze PID %d: %w", err)
	}

	return nil
}

// ThawPID unfreezes the process and removes it from the isolation group.
func ThawPID(pid int) error {
	groupPath := filepath.Join(CGROUP_BASE, fmt.Sprintf("proc_%d", pid))
	
	// 1. Trigger the thaw
	freezeFile := filepath.Join(groupPath, "cgroup.freeze")
	if err := os.WriteFile(freezeFile, []byte("0"), 0644); err != nil {
		return fmt.Errorf("failed to thaw PID %d: %w", err)
	}

	// 2. Cleanup: Move process back to root cgroup
	rootProcs := "/sys/fs/cgroup/cgroup.procs"
	if err := os.WriteFile(rootProcs, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		fmt.Printf("[Warden Warning] Failed to return PID %d to root cgroup: %v\n", pid, err)
	}

	// 3. Delete the group directory
	if err := os.Remove(groupPath); err != nil {
		return fmt.Errorf("failed to cleanup cgroup directory: %w", err)
	}

	return nil
}
