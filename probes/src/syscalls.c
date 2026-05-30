/*
 * REPOSITORY: PhoenixKernel
 * ARCHITECTURAL JUSTIFICATION: Structured Syscall Boundary for per-process monitoring.
 * DEPENDENCY BOUNDARY: BPF side only.
 * DETERMINISTIC CONSIDERATIONS: Fast-path bitmask lookups for O(1) validation.
 */

#include "../src/vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_core_read.h>

// Max number of syscalls to track in the bitmask (Linux has ~450)
#define MAX_SYSCALLS 512
#define BITMASK_WORDS (MAX_SYSCALLS / 64)

struct syscall_event {
    u32 pid;
    u32 syscall_nr;
    u32 entropy_flag; // 0 = STABLE, 1 = HIGH_ENTROPY
    u64 timestamp;
    u8 comm[16];
};

struct profile {
    u64 bitmask[BITMASK_WORDS];
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 256 * 1024);
} rb_syscalls SEC(".maps");

// PID -> Syscall Bitmask Profile
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 2048);
    __type(key, u32);
    __type(value, struct profile);
} syscall_profiles SEC(".maps");

SEC("raw_tracepoint/sys_enter")
int handle_syscall_entry(struct bpf_raw_tracepoint_args *ctx) {
    u32 pid = bpf_get_current_pid_tgid() >> 32;
    u32 syscall_nr = (u32)ctx->args[1];
    struct profile *prof;
    struct syscall_event *e;

    prof = bpf_map_lookup_elem(&syscall_profiles, &pid);
    
    // If no profile exists, or the syscall bit is not set, it's HIGH_ENTROPY
    u32 entropy = 1;
    if (prof) {
        u32 word_idx = syscall_nr / 64;
        u32 bit_idx = syscall_nr % 64;
        if (word_idx < BITMASK_WORDS) {
            if (prof->bitmask[word_idx] & (1ULL << bit_idx)) {
                entropy = 0; // STABLE
            }
        }
    }

    // Optimization: Only report HIGH_ENTROPY events for now to prevent bus saturation
    // In a full implementation, STABLE events would be sampled.
    if (entropy == 1) {
        e = bpf_ringbuf_reserve(&rb_syscalls, sizeof(*e), 0);
        if (!e) return 0;

        e->pid = pid;
        e->syscall_nr = syscall_nr;
        e->entropy_flag = entropy;
        e->timestamp = bpf_ktime_get_ns();
        bpf_get_current_comm(&e->comm, sizeof(e->comm));

        bpf_ringbuf_submit(e, 0);
    }

    return 0;
}

char LICENSE[] SEC("license") = "GPL";
