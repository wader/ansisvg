package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wader/ansisvg/ansisvg"
	"github.com/wader/ansisvg/colorscheme/schemes"
)

type dimension struct {
	Width  int
	Height int
}

func (d *dimension) String() string {
	return fmt.Sprintf("%dx%d", d.Width, d.Height)
}

func (d *dimension) Set(s string) error {
	parts := strings.Split(s, "x")
	if len(parts) != 2 {
		return fmt.Errorf("must be WxH")
	}
	d.Width, _ = strconv.Atoi(parts[0])
	d.Height, _ = strconv.Atoi(parts[1])
	return nil
}

var fontFlag = flag.String("font", "Monaco, Lucida Console, Courier", "Font")
var fontSizeFlag = flag.Int("fontsize", 12, "Font size")
var terminalWidthFlag = flag.Int("width", 0, "Terminal width (auto)")
var characterBox = dimension{Width: 7, Height: 13}
var colorSchemeFlag = flag.String("colorscheme", "Builtin Dark", "Color scheme")

func init() {
	flag.Var(&characterBox, "chardimension", "Character box dimension")
}

func run() error {
	flag.Parse()

	colorScheme, err := schemes.Load(*colorSchemeFlag)
	if err != nil {
		return err
	}

	return ansisvg.Convert(
		os.Stdin,
		os.Stdout,
		ansisvg.Options{
			Font:          *fontFlag,
			FontSize:      *fontSizeFlag,
			TerminalWidth: *terminalWidthFlag,
			CharacterBox: ansisvg.Dimension{
				Width:  characterBox.Width,
				Height: characterBox.Height,
			},
			ColorScheme: colorScheme,
		},
	)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
