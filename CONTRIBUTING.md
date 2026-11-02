# Contributing to warmbench

Thanks for your interest. warmbench is at the design stage: the documents in `docs/` are the
source of truth, and no feature works yet.

## Before v0.1

- **Open an issue first** for anything beyond a typo fix, so the change can be checked against
  [docs/DESIGN.md](docs/DESIGN.md) and [docs/POOL.md](docs/POOL.md).
- **Measurements beat opinions.** Claims about speed, memory or disk need numbers from a real
  Unity project, with the Unity version, machine and project size stated. Add them to
  [docs/spike.md](docs/spike.md) or attach them to the issue.
- **Design changes** go in the design documents in the same pull request as the code.

## Project structure

```
cmd/warmbench/          CLI client
cmd/warmbenchd/         daemon
internal/               scheduler, slot, store, governor, driver
unity/com.warmbench.agent/   editor-only Unity package (C#)
docs/                   DESIGN.md, POOL.md, spike.md, diagrams/
```

Keep the scheduler, store and MCP layer free of Unity types; Unity-specific code lives in
`internal/driver` and `unity/`.

## Development

- Go 1.23 or later: `make build`, `make test`, `make vet`, `make fmt-check`.
- Unity 6000.x (or 2022.3) to work on the editor package. Never commit a `Library`, `Temp` or
  `Logs` folder.

## Pull requests

- One topic per pull request, with a description of what changed and how it was tested.
- CI must pass (`gofmt`, `go vet`, `go build`, `go test`).
- Add a line under `Unreleased` in [CHANGELOG.md](CHANGELOG.md) for user-visible changes.

## License

By contributing you agree that your contributions are licensed under the MIT License.
