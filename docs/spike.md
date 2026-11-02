# Spike: measure before building

Milestone 1 (DESIGN.md section 10). Nothing below has been measured yet. The go / no-go decision
for v0.1 and for each memory mechanism in POOL.md section 10 depends on these numbers.

## Setup

- A copy of a real Unity project outside its own repository (the slot project).
- One editor kept open with `-executeMethod Warmbench.Agent.Serve`, a minimal socket listener and
  a `TestRunnerApi` run for a filter.
- Record: Unity version, OS and kernel, CPU, RAM, disk type, filesystem, project size
  (`Assets/` and `Library/`), test counts.
- Do not run while other Unity editors on the machine are in use.

## Speed

| Case | Cold run | Warm (HOT) run | Notes |
|---|---|---|---|
| no change | | | |
| one script change | | | |
| one asset change | | | |
| package manifest change | | | |

## Warmth tiers (POOL.md 10.2, 10.6)

| Measure | Result |
|---|---|
| warm boot from an existing `Library` (WARM) | |
| thaw of a frozen, reclaimed editor, no pre-thaw | |
| thaw with pre-thaw | |
| resident memory of a frozen, reclaimed editor | |
| licence session still valid after a long freeze | |

## Disk and memory (POOL.md 6, 10.3, 10.4, 10.8)

| Measure | Result |
|---|---|
| reflink clone time of a seed | |
| does a cloned `Library` boot without reimport | |
| overlay seed: copy-up cost, import speed, page-cache saving with 2 to 4 slots | |
| KSM: memory merged across 2 to 4 editors, `ksmd` CPU cost | |
| baseline and peak RSS per class | |
| RSS growth per job after in-editor cleanup | |
| per-slot `HOME` accepted by the licence check | |

## Other

| Measure | Result |
|---|---|
| shadow compile agreement with Unity's compiler | |
| share of agent iterations that fail at compile | |
| domain reloads before a slot needs recycling | |
