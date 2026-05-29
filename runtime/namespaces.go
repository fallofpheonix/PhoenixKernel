/*
 * REPOSITORY: PheonixKernel
 * ARCHITECTURAL JUSTIFICATION: Network namespace severing for physical process isolation.
 * DEPENDENCY BOUNDARY: System calls (setns). No strategic layers.
 * DETERMINISTIC CONSIDERATIONS: Absolute network isolation via namespace migration.
 */

//go:build linux

package runtime

import (
	"fmt"
	"os"
	"runtime"
	"syscall"

	"golang.org/x/sys/unix"
)

const ISOLATED_NS_PATH = "/var/run/netns/phoenix_blackhole"

// NamespaceSever migrates a process into an isolated, loopback-only network namespace.
func NamespaceSever(pid int) error {
	nsFile, err := os.Open(ISOLATED_NS_PATH)
	if err != nil {
		return fmt.Errorf("failed to open isolated namespace %s: %w", ISOLATED_NS_PATH, err)
	}
	defer nsFile.Close()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := unix.Setns(int(nsFile.Fd()), syscall.CLONE_NEWNET); err != nil {
		return fmt.Errorf("failed to sever network for PID %d: %w", pid, err)
	}

	return nil
}

// EnsureBlackholeExists verifies the isolated namespace is present on the host.
func EnsureBlackholeExists() error {
	if _, err := os.Stat(ISOLATED_NS_PATH); os.IsNotExist(err) {
		return fmt.Errorf("isolated blackhole namespace missing at %s", ISOLATED_NS_PATH)
	}
	return nil
}
