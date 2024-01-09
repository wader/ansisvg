package main

import (
	"fmt"
	"os"

	"github.com/wader/ansisvg/cli"
)

// set by release build
// -ldflags "-X main.version=1.2.3"
var version string = "dev"

func main() {
	if err := cli.Main(cli.Env{
		Version:  version,
		ReadFile: os.ReadFile,
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		Args:     os.Args,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
