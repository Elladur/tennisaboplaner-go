# AGENTS.md

Guide for AI coding agents working in `tennisaboplaner-go`.

## Project

Go CLI that schedules tennis matches over a season. Reads `settings.json`, optimizes a schedule (equal play time, even matchups, even distribution), and exports an Excel file plus per-player `.ics` calendars into `output/`.

Module: `github.com/Elladur/tennisaboplaner-go`. Go 1.27.

## Commands

- Build: `go build ./...`
- Run: `go run . --config settings.json --outDir output --times 100`
- Tests: `go test ./...`
- Tests (verbose): `go test -v ./internal/...`
- Integration test (performance): `go test ./test/integration/...`
- Integration tests (flaky/tagged): `go test -tags=integration ./test/integration/...`
- Format: `gofmt -w .`
- Vet: `go vet ./...`
- Lint: `golangci-lint run ./...`

After every change, run the full check suite: `gofmt -w .`, `go vet ./...`, `golangci-lint run ./...`, and `go test ./...` before considering the work done.

CLI flags (defined in `cmd/root.go`): `--config` (default `settings.json`), `--level` (`info`/`warn`/`debug`), `--outDir` (default `output`), `--times` (default `100`).

## Structure

- `main.go` — entrypoint, calls `cmd.Execute()`.
- `cmd/` — Cobra root command; flag parsing, settings load, calls `internal.ExecutePlanerParallel`.
- `internal/` — domain logic:
  - `season.go` — `Season`, `SeasonSettings`, `CreateSeasonFromSettings`.
  - `player.go` — `Player` (JSON marshal/unmarshal).
  - `match.go` — `Match` (JSON marshal/unmarshal, `createMatch`).
  - `schedule.go` — schedule helpers.
  - `round.go` — round helpers.
  - `optimizer.go` — `Optimizer.Optimize()` (swaps players/matches until no improvement).
  - `scoring_algorithm.go` — `GetScore` and std-dev metrics.
  - `executor.go` — serial/parallel execution of planning runs, picks best score.
  - `exporter.go` — Excel + iCal export.
  - `helper.go` — generic helpers (`isInSlice`, `pop`, `shuffle`).
- `test/integration/` — performance and tagged integration tests; testdata lives in `test/integration/testdata/`.
- `test/cpu_perfm.go` — standalone perf profiling helper.
- `settings.json` / `setting_2027.json` — example configs.

## Conventions

- All domain code lives in `internal/`; only `cmd` and `main` are outside.
- Tests are co-located in `internal/*_test.go` and use testify + `assert`.
- Shared test fixtures: `setupTestSeason()`, `setupStaticTestSeason()` in `season_test.go`.
- Keep player indices as `int` (0-based) in `Match`/schedule; `Player` slice index == player id.
- `Match.player1` is always the lower ID; use `createMatch` to enforce ordering.
- Mutating fields on `Season` (e.g. `dates`, `partialRounds`) are unexported; construct via `CreateSeasonFromSettings`/`createSeason`.
- Logging via logrus; level controlled by `--level`.
- Output dir and `output/` are gitignored; don't commit generated files.
- Add unit tests and integration tests for new or changed behavior whenever feasible; keep them co-located with the code (`internal/*_test.go`, `test/integration/`).

## Notes

- `go.mod` declares `go 1.27.3` (toolchain). Use the installed Go.
- Integration test reads `test/integration/input/performance.json` with relative path; run from that directory's context or repo root as `go test ./test/integration/...` (it resolves relative to package dir).
- Parallel execution uses goroutines + channel; `ExecutePlanerSerial` exists for deterministic/perf testing.