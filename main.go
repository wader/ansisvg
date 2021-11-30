package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wader/ansisvg/ansitosvg"
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

var fontNameFlag = flag.String("fontname", ansitosvg.DefaultOptions.FontName, "Font name")
var fontSizeFlag = flag.Int("fontsize", ansitosvg.DefaultOptions.FontSize, "Font size")
var terminalWidthFlag = flag.Int("width", 0, "Terminal width (auto)")
var characterBoxSize = boxSize{
	Width:  ansitosvg.DefaultOptions.CharacterBoxSize.Width,
	Height: ansitosvg.DefaultOptions.CharacterBoxSize.Height,
}
var colorSchemeFlag = flag.String("colorscheme", ansitosvg.DefaultOptions.ColorScheme, "Color scheme")
var transparentFlag = flag.Bool("transparent", ansitosvg.DefaultOptions.Transparent, "Transparent background")

func init() {
	flag.Var(&characterBoxSize, "charboxsize", "Character box size")
}

func main() {
	flag.Parse()

	if err := ansitosvg.Convert(
		os.Stdin,
		os.Stdout,
		ansitosvg.Options{
			FontName:      *fontNameFlag,
			FontSize:      *fontSizeFlag,
			TerminalWidth: *terminalWidthFlag,
			CharacterBoxSize: ansitosvg.BoxSize{
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
