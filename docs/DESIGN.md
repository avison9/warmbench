# warmbench: design

Status: draft, 2026-11-02. Nothing is built yet.

## 1. The problem

A team (or one developer) runs several AI coding agents in parallel on one Unity project, each in its own git worktree. Every agent must run the Unity Test Framework (edit mode and play mode) before it pushes. The machine has one Unity install and, usually, one licence seat.

What happens today, observed on the bleed project (Unity 6000.3, four agents, one Linux workstation, 2026-11-02):

- Each test run is a cold start: launch the editor in `-batchmode`, open the worktree's project, import or refresh its `Library`, compile every assembly, reload the domain, and only then run tests. A full play-mode run took 15 to 20 minutes; most of that is not test time.
- Each worktree carries its own `Library` copy (several GB) that has to be copied in and later deleted.
- Running editors side by side fights for RAM and the GPU, and timing-sensitive play-mode tests then fail for reasons that are not the code. So the agents serialise with `flock` on a lock file.
- The lock is blind. An agent cannot see its place in the queue, an estimate of the wait, or who holds the lock. There is no priority (a PR's final pre-merge check waits behind an exploratory run), and identical trees are tested twice.
- Four agents with review rounds each spent most of the afternoon waiting in that queue.

Unity itself refuses to open the same project path twice ("only a single instance of Unity can run at a time" for one project in batch mode). Separate worktrees are separate paths, so the limit is cost and contention, not a hard rule.

## 2. What exists, and the gap

| Tool | Helps with | Does not help with |
|---|---|---|
| Unity Accelerator (cache server) | shares imported asset artifacts across `Library` copies | cold start per run; no queue |
| GameCI (Docker images, GitHub Actions) | reproducible CI runs | built for CI, not interactive agents; still cold |
| Unity Build Automation | cloud builds | paid; not a local test loop |
| ParrelSync, Multiplayer Play Mode | several editors on one project via links | made for multiplayer testing, not scheduling |
| Unity MCP servers (CoplayDev, StraySpark, AgentBridge, mcp-unity) | an agent drives an open editor, including `run_tests` | one editor per project; no scheduler across agents (AgentBridge has a file queue for a single editor) |

Gap: nothing owns the machine's editors as a shared, warm, scheduled resource for many agents.

## 3. Goals and non-goals

Goals
1. A test run on a small change finishes in minutes, not tens of minutes, by reusing an editor that is already open with a warm `Library`.
2. Agents queue fairly and can see their position, the holder and an ETA.
3. Identical inputs are never tested twice.
4. Results are structured (pass/fail per test, failing tests' logs only), ready for an agent to read.
5. Disk is cleaned up by the service, not by each agent.
6. Works on one machine with one licence by default.

Non-goals (first versions)
- Builds (APK, AAB, IPA). Possible later as another job type.
- Replacing CI. It can serve CI too, but the target is the local agent loop.
- Unreal or Godot. The design should not prevent them (see 11).

## 4. Architecture

The pool design below is superseded by [POOL.md](POOL.md) (seeds, build-target classes, verification ladder, cost-based placement). This section is kept as the baseline it builds on.

```
 agents (Claude Code, Codex, scripts, humans)
     |  MCP (stdio or HTTP) / CLI / HTTP API
     v
 +-----------------------------+
 | warmbenchd (Go daemon)      |
 |  - job queue + scheduler    |
 |  - result cache (by hash)   |
 |  - slot manager             |
 |  - resource monitor         |
 +-----------------------------+
     | local socket per slot (JSON lines)
     v
 slot 1: Unity editor (batchmode, stays up)   slot 2: ...
   project copy: /var/lib/warmbench/slots/1/<project>
   warmbench editor package (C#) listening for jobs
```

### 4.1 Slots

A slot is one long-lived Unity editor on its own copy of the project. The copy starts from the base branch (for example `origin/main`). The editor is launched once:

```
Unity -batchmode -nographics? -projectPath <slot copy> -executeMethod Warmbench.Agent.Serve -logFile <slot log>
```

(Play-mode tests that need graphics run in slots launched without `-nographics`; the slot records which kind it is.)

`Warmbench.Agent.Serve` keeps the editor alive and listens on a Unix socket (named pipe on Windows). It never calls `EditorApplication.Exit` unless asked.

### 4.2 A job's life

1. **Submit.** The agent submits: worktree path (or a git ref), test platform (EditMode / PlayMode), filter (names, categories, assemblies), priority, and a label.
2. **Key.** The daemon computes a content key: the git tree hash of the worktree's `Assets/`, `Packages/`, `ProjectSettings/` (including uncommitted changes, via `git stash create` or by hashing the working tree) plus Unity version, platform, filter and defines. If the cache has this key, return the stored result at once.
3. **Schedule.** Pick a free slot whose kind fits (graphics or not). Prefer the slot whose current tree is closest to the job's tree (fewest changed files), so the least is recompiled and reimported.
4. **Sync.** Make the slot's copy equal to the job's tree. Use `rsync --delete --checksum` over the three folders, or a git checkout in the slot copy followed by applying the worktree's uncommitted diff. Never touch the slot's `Library`.
5. **Refresh.** Tell the editor to `AssetDatabase.Refresh()` and wait for compilation and domain reload to finish (`CompilationPipeline` events, `EditorApplication.isCompiling`). Compile errors end the job here with the errors as the result.
6. **Run.** Use the Test Framework's `TestRunnerApi` with the job's filter. Stream progress (tests done / total) to the daemon.
7. **Report.** Return a compact result: counts, failures with message and stack, the slot log excerpt for failing tests only, durations. Store the full NUnit XML and log for later.
8. **Reset if needed.** If the editor crashed, hung (no progress for N minutes), or ended in a bad state (play mode stuck, domain not reloaded), kill and restart the slot and requeue the job once.

### 4.3 Expected gain

Cold run today: start + open + import + compile + reload + tests. Warm run: sync + incremental import of changed assets + compile of changed assemblies + reload + tests. For a change that touches a few scripts, the guess is 1 to 3 minutes of overhead instead of 10 or more. This must be measured (see 10) before any claim is published.

## 5. Scheduler

- **Priority classes:** `final` (pre-merge full suite), `normal`, `explore`. Within a class, first come first served, with ageing so nothing starves.
- **Deduplication:** a job whose key matches a running job attaches to it instead of queueing.
- **Change-aware selection (optional, later):** map changed files to assemblies (asmdef graph) and offer "affected tests only" for early runs; the `final` class always runs the full suite.
- **Sharding (later):** split one play-mode suite across several slots when the machine allows, merge the results.
- **ETA:** keep a moving average of duration per (platform, filter) and per test; estimate queue wait from the jobs ahead.
- **Cancel and supersede:** a new job from the same worktree and label can cancel the older one still queued.

## 6. Resources

- Number of slots from config, capped by measured RAM (an editor plus play mode on a mid-size project is several GB) and CPU cores. Default: 1 graphics slot + 1 headless slot.
- Run graphics-dependent play-mode tests only in graphics slots, one at a time, to keep timing tests stable.
- Watch memory; if a slot's editor grows past a limit (editors leak over many domain reloads), recycle it between jobs.

## 7. Interfaces

### 7.1 MCP server (for agents)

Tools:
- `run_tests(worktree, platform, filter?, priority?, label?, wait?)` returns a job id, or the result if `wait` and it finishes in time.
- `job_status(job_id)` returns state (queued with position and ETA, syncing, compiling, running with progress, done), and the result when done.
- `queue()` returns all jobs and slots, who is running, and waits.
- `cancel(job_id)`.
- `logs(job_id, test?)` returns the log excerpt for one failing test.

### 7.2 CLI

`warmbench run --platform playmode --filter Gallery --wait`, `warmbench status`, `warmbench queue`, `warmbench cancel <id>`, `warmbench slots`, `warmbench gc`. The CLI is a thin client of the daemon's HTTP API on localhost. Exit code follows the test result, so it can replace `tools/unity/test.sh` in existing scripts.

### 7.3 Editor package

A Unity package (`com.warmbench.agent`, editor-only assembly) with `Serve`, the socket loop, and the refresh / run / report logic. No runtime code ships in player builds.

## 8. Disk and cleanup

- Slot copies are the only `Library` copies; agents' worktrees do not need one at all.
- `warmbench gc` removes stored results and logs older than a set age, and slot copies for projects no longer configured.
- The daemon reports disk use per slot.

## 9. Licensing and safety

- Default: one machine, one seat; the daemon never starts more editors than configured. Whether several editors on one machine are within a given Unity licence must be checked against Unity's current terms; the README must say so plainly and not promise it.
- The daemon runs as the user, listens on localhost or a Unix socket only, and executes only the Test Framework through its own editor package: it does not run arbitrary methods from requests.
- A job's worktree path must be inside configured roots.

## 10. Milestones

1. **Spike (measure first).** On bleed: one slot, manual sync, `TestRunnerApi` run through a socket. Measure cold run against warm run for (a) a one-script change, (b) an asset change, (c) no change. Go or no-go on these numbers.
2. **v0.1.** Daemon with queue, one or two slots, rsync sync, CLI, result cache by key, crash recovery. Replace bleed's `flock` wrapper with `warmbench run`.
3. **v0.2.** MCP server; ETA; priorities; dedupe; gc.
4. **v0.3.** Change-aware selection; sharding; dashboard (simple web page from the daemon).
5. **Later.** Builds as jobs; remote slots on other machines for teams with more seats; other engines.

## 11. Open questions

- Sync method: rsync of three folders vs git worktree per slot plus diff apply. rsync is simpler and handles uncommitted files; measure its time on a large `Assets/`.
- Domain reload reliability over hundreds of jobs in one editor; how often must a slot be recycled?
- Package changes (`manifest.json`) force a package resolve; treat them as "recycle slot" jobs?
- Unity version per project: one daemon serving projects on different editor versions needs slots per version.
- Windows and macOS: named pipes vs Unix sockets; the Go daemon should abstract this from day one.
- Engine-agnostic core: keep the scheduler, cache and MCP layer free of Unity types so an Unreal or Godot "slot driver" can be added.
- Name: `warmbench` is a working name; check it is free on GitHub, pkg.go.dev and openupm.

## 12. Prior art and references

- Unity command line arguments: https://docs.unity3d.com/Manual/CommandLineArguments.html
- Unity Test Framework `TestRunnerApi`: https://docs.unity3d.com/Packages/com.unity.test-framework@latest
- Unity Accelerator: https://docs.unity3d.com/Manual/UnityAccelerator.html
- GameCI: https://game.ci
- AgentBridge (file queue for one editor): https://glama.ai/mcp/servers/simonwittber/AgentBridge
- StraySpark Unity MCP server (`run_tests`): https://www.strayspark.studio/docs/unity-mcp-server
- mcp-unity (CoderGamester): https://archestra.ai/mcp-catalog/codergamester__mcp-unity
- Unity forum on build machine licensing: https://discussions.unity.com/t/unity-licensing-for-build-machines/608665
