// Package main is the APM CLI Go entry point.
// This is a stub that will grow as more Python modules are migrated.
package main

import (
	"fmt"
	"os"

	"github.com/githubnext/apm/internal/version"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(version.GetVersion())
		return
	}
	fmt.Fprintln(os.Stderr, "apm-go: stub binary (migration in progress)")
	os.Exit(1)
}
