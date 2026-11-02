// Package scheduler orders jobs and places them on slots.
//
// It owns the queue (priority classes with ageing and per-agent fair share),
// dedupe of identical job keys, the verification ladder planner, and
// cost-based placement of a job on the slot closest to its tree.
// See docs/POOL.md sections 7 and 8.
package scheduler
