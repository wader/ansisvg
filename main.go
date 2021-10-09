package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wader/ansisvg/ansisvg"
)

type boxSize struct {
	Width  int
	Height int
}

func (d *boxSize) String() string {
	return fmt.Sprintf("%dx%d", d.Width, d.Height)
}

func (d *boxSize) Set(s string) error {
	parts := strings.Split(s, "x")
	if len(parts) != 2 {
		return fmt.Errorf("must be WxH")
	}
	d.Width, _ = strconv.Atoi(parts[0])
	d.Height, _ = strconv.Atoi(parts[1])
	return nil
}

var fontFlag = flag.String("font", ansisvg.DefaultOptions.Font, "Font")
var fontSizeFlag = flag.Int("fontsize", ansisvg.DefaultOptions.FontSize, "Font size")
var terminalWidthFlag = flag.Int("width", 0, "Terminal width (auto)")
var characterBoxSize = boxSize{
	Width:  ansisvg.DefaultOptions.CharacterBoxSize.Width,
	Height: ansisvg.DefaultOptions.CharacterBoxSize.Height,
}
var colorSchemeFlag = flag.String("colorscheme", ansisvg.DefaultOptions.ColorScheme, "Color scheme")

func init() {
	flag.Var(&characterBoxSize, "charboxsize", "Character box size")
}

func run() error {
	flag.Parse()

	return ansisvg.Convert(
		os.Stdin,
		os.Stdout,
		ansisvg.Options{
			Font:          *fontFlag,
			FontSize:      *fontSizeFlag,
			TerminalWidth: *terminalWidthFlag,
			CharacterBoxSize: ansisvg.BoxSize{
				Width:  characterBoxSize.Width,
				Height: characterBoxSize.Height,
			},
			ColorScheme: *colorSchemeFlag,
		},
	)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
