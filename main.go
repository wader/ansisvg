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

var fontNameFlag = flag.String("fontname", ansisvg.DefaultOptions.FontName, "Font name")
var fontSizeFlag = flag.Int("fontsize", ansisvg.DefaultOptions.FontSize, "Font size")
var terminalWidthFlag = flag.Int("width", 0, "Terminal width (auto)")
var characterBoxSize = boxSize{
	Width:  ansisvg.DefaultOptions.CharacterBoxSize.Width,
	Height: ansisvg.DefaultOptions.CharacterBoxSize.Height,
}
var colorSchemeFlag = flag.String("colorscheme", ansisvg.DefaultOptions.ColorScheme, "Color scheme")
var transparentFlag = flag.Bool("transparent", ansisvg.DefaultOptions.Transparent, "Transparent background")

func init() {
	flag.Var(&characterBoxSize, "charboxsize", "Character box size")
}

func main() {
	flag.Parse()

	if err := ansisvg.Convert(
		os.Stdin,
		os.Stdout,
		ansisvg.Options{
			FontName:      *fontNameFlag,
			FontSize:      *fontSizeFlag,
			TerminalWidth: *terminalWidthFlag,
			CharacterBoxSize: ansisvg.BoxSize{
				Width:  characterBoxSize.Width,
				Height: characterBoxSize.Height,
			},
			ColorScheme: *colorSchemeFlag,
			Transparent: *transparentFlag,
		},
	); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
