# Changelog

All notable changes to warmbench are recorded here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Design documents: `docs/DESIGN.md` (problem, goals, interfaces, milestones) and `docs/POOL.md`
  (pool architecture: seeds, lanes, verification ladder, cost-based placement, memory and disk
  discipline, slot lifecycle, results for agents).
- Diagrams in `docs/diagrams/`.
- Milestone 1 measurement plan in `docs/spike.md`.
- Project scaffold: Go module with `warmbench` and `warmbenchd` stubs, internal package layout,
  the `com.warmbench.agent` Unity editor package stub, and CI.
