/*
 * REPOSITORY: PhoenixKernel
 * ARCHITECTURAL JUSTIFICATION: eBPF sensor harness for network events.
 * DEPENDENCY BOUNDARY: BPF side only.
 * DETERMINISTIC CONSIDERATIONS: Zero kernel latency degradation.
 */

#include "../src/vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_core_read.h>

struct net_event {
    u32 pid;
    u32 uid;
    u16 family;
    u16 sport;
    u16 dport;
    u32 saddr;
    u32 daddr;
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 256 * 1024);
} rb_net SEC(".maps");

SEC("tp/syscalls/sys_enter_connect")
int handle_connect(struct trace_event_raw_sys_enter *ctx) {
    struct net_event *e;
    struct sockaddr_in *addr;

    e = bpf_ringbuf_reserve(&rb_net, sizeof(*e), 0);
    if (!e) return 0;

    e->pid = bpf_get_current_pid_tgid() >> 32;
    e->uid = bpf_get_current_uid_gid();

    addr = (struct sockaddr_in *)ctx->args[1];
    bpf_probe_read_user(&e->family, sizeof(e->family), &addr->sin_family);
    bpf_probe_read_user(&e->dport, sizeof(e->dport), &addr->sin_port);
    bpf_probe_read_user(&e->daddr, sizeof(e->daddr), &addr->sin_addr.s_addr);

    bpf_ringbuf_submit(e, 0);
    return 0;
}

char LICENSE[] SEC("license") = "GPL";
