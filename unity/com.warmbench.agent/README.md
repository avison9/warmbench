# com.warmbench.agent

Editor-only Unity package used by warmbench slots. A slot editor is launched with

```
Unity -batchmode -projectPath <slot> -executeMethod Warmbench.Agent.Serve -logFile <log>
```

`Serve` will keep the editor alive, listen on a local socket for jobs from warmbenchd, refresh
and compile, run `TestRunnerApi` with the job's filter, stream progress and report results, then
clean up. See `docs/DESIGN.md` section 4 and `docs/POOL.md` in the repository root.

Status: scaffold only. No runtime code ships in player builds.
