// Command warmbench is the command-line client of warmbenchd.
//
// It is a thin client of the daemon's HTTP API on localhost. Nothing is
// implemented yet; see docs/DESIGN.md and docs/POOL.md.
package main

import (
	"fmt"
	"os"
)

const usage = `warmbench: run Unity tests on warm editors shared by many coding agents

Status: design and scaffold only. No command works yet.
See docs/DESIGN.md and docs/POOL.md.

Usage:
  warmbench run --platform <editmode|playmode> [--filter <filter>] [--wait]
  warmbench status <job-id>
  warmbench queue
  warmbench cancel <job-id>
  warmbench slots
  warmbench gc
  warmbench doctor
`

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Print(usage)
		return
	}
	fmt.Fprintf(os.Stderr, "warmbench %s: not implemented yet\n", os.Args[1])
	os.Exit(2)
}
