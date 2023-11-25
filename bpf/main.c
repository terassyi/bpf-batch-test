#include "vmlinux.h"
#include <bpf/bpf_helpers.h>

struct
{
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(key_size, sizeof(u32));
	__uint(value_size, sizeof(u32));
	__uint(max_entries, 1024 * 1024 * 100);
} test_map SEC(".maps");

SEC("xdp")
int nop(struct xdp_md *ctx)
{
	return XDP_PASS;
}

char LICENSE[] SEC("license") = "Dual BSD/GPL";
