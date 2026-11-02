// Package governor hands out resource tokens (ram, cpu, gpu, quiet, device)
// and enforces memory limits: computed max_hot, reservations by p95 peak,
// pressure (PSI) handling and per-slot cgroup limits on Linux.
// See docs/POOL.md sections 9 and 10.
package governor
