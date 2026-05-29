//go:build !linux

package runtime

import "fmt"

// NamespaceSever is a stub for non-linux systems.
func NamespaceSever(pid int) error {
	fmt.Printf("[SIMULATION] NamespaceSever called for PID %d on non-linux host\n", pid)
	return nil
}

// EnsureBlackholeExists is a stub for non-linux systems.
func EnsureBlackholeExists() error {
	return nil
}
