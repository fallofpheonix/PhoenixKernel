package live

// ProbeLoader dynamically attaches eBPF objects to the kernel
type ProbeLoader struct{}

func (p *ProbeLoader) Load() error { return nil }
