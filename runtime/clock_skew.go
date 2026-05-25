package runtime

// ClockSkew monitors and corrects drift between core clocks.
type ClockSkew struct {
	DriftNs int64
}

func (cs *ClockSkew) RecordSkew(coreA, coreB int64) {
	cs.DriftNs = coreA - coreB
}

func (cs *ClockSkew) GetDrift() int64 {
	return cs.DriftNs
}
