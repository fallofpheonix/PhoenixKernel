package probes

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang bpf src/syscalls.c -- -I../src
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang netbpf src/network.c -- -I../src
