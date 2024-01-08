package main

import (
	"fmt"
	"os"

	"github.com/wader/ansisvg/cli"
)

func main() {
	if err := cli.Main(cli.Env{
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
