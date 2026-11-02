// Command warmbenchd is the warmbench daemon: job queue and scheduler,
// result cache, seed and slot manager, and resource governor.
//
// Nothing is implemented yet; see docs/POOL.md for the architecture.
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "warmbenchd: not implemented yet; see docs/POOL.md")
	os.Exit(2)
}
