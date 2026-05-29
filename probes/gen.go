package probes

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang bpf syscalls.c -- -I../src
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang netbpf network.c -- -I../src
