# REPO_OVERVIEW

Purpose: PheonixKernel contains low-level kernel integration code (eBPF, loaders, and helpers) used by the PhoenixOS substrate to perform kernel-space monitoring and reflexive actuation. This repo provides the eBPF object artifacts and the Go-based loader/adapter used at runtime.

Primary source documents used to derive this overview:
- `go.mod` — shows `github.com/cilium/ebpf` dependency and indicates eBPF usage.
- `src/` and `PheonixKernel/src/phoenix_exec.o` — eBPF object used by runtime images.

Suggested next steps: add a README.md with build instructions for the eBPF object and how to run loader tests.
