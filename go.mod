module github.com/fallofpheonix/PheonixKernel

go 1.26

replace github.com/fallofpheonix/PheonixCore => ../PheonixCore

require (
	github.com/cilium/ebpf v0.17.3
	github.com/fallofpheonix/PheonixCore v0.0.0
	golang.org/x/sys v0.30.0
)
