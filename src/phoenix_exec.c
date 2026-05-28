#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_core_read.h>

struct exec_event {
    u32 pid;
    u32 ppid;
    u32 tgid; // Thread Group ID
    u32 nsproxy_ino; // Namespace proxy inode
    u32 uid;
    u8 comm[16];
    u8 filename[128];
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 256 * 1024);
} rb SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1024);
    __type(key, u32);
    __type(value, u32);
} blocked_pids SEC(".maps");

SEC("tp/syscalls/sys_enter_execve")
int handle_execve(struct trace_event_raw_sys_enter *ctx) {
    struct exec_event *e;
    struct task_struct *task;
    struct ns_common ns = {};
    struct mnt_namespace *mnt_ns;

    e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
    if (!e) return 0;

    task = (struct task_struct *)bpf_get_current_task();
    e->pid = bpf_get_current_pid_tgid() >> 32;
    e->tgid = (u32)bpf_get_current_pid_tgid();
    e->ppid = BPF_CORE_READ(task, real_parent, tgid);
    mnt_ns = BPF_CORE_READ(task, nsproxy, mnt_ns);
    if (mnt_ns != NULL) {
        bpf_core_read(&ns, sizeof(ns), &mnt_ns->ns);
        e->nsproxy_ino = ns.inum;
    }
    e->uid = bpf_get_current_uid_gid();
    
    bpf_get_current_comm(&e->comm, sizeof(e->comm));
    
    // Read filename from the first syscall argument
    bpf_probe_read_user_str(&e->filename, sizeof(e->filename), (void *)ctx->args[0]);

    bpf_ringbuf_submit(e, 0);
    return 0;
}

SEC("lsm/bprm_check_security")
int BPF_PROG(phoenix_enforce_exec, struct linux_binprm *bprm) {
    u32 pid = bpf_get_current_pid_tgid() >> 32;
    u32 *action;

    action = bpf_map_lookup_elem(&blocked_pids, &pid);
    if (action && *action == 1) {
        // Return -1 (Operation not permitted) to block the execve syscall
        return -1;
    }

    return 0;
}

SEC("lsm/file_mprotect")
int BPF_PROG(phoenix_enforce_mprotect, struct vm_area_struct *vma, unsigned long reqprot, unsigned long prot) {
    u32 pid = bpf_get_current_pid_tgid() >> 32;
    u32 *action;

    action = bpf_map_lookup_elem(&blocked_pids, &pid);
    if (action && *action == 1) {
        // Return -1 (Operation not permitted) to block the mprotect syscall
        return -1;
    }

    return 0;
}

char LICENSE[] SEC("license") = "GPL";
