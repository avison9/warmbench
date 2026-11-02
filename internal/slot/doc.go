// Package slot manages the lifecycle of slots: long-lived editor processes
// grown from seeds.
//
// It owns the slot state machine, leases with epoch fencing, heartbeats and
// the watchdog, warmth tiers (hot, frozen, warm), recycling and per-job
// cleanup. See docs/POOL.md sections 10 and 11.
package slot
