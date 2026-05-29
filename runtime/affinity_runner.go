package runtime


// AffinityRunner manages CPU affinity for deterministic replay.
type AffinityRunner struct {
	CoreID int
}

func (ar *AffinityRunner) LockToCore(core int) error {
	// Note: Actual affinity locking requires syscalls like sched_setaffinity
	// which varies by OS. This is a platform-agnostic scaffold.
	ar.CoreID = core
	return nil
}

func (ar *AffinityRunner) CurrentCore() int {
	return ar.CoreID
}
