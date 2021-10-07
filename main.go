package main

import (
	"github.com/wader/ansisvg/ansidecoder"
	"github.com/wader/ansisvg/colorscheme/schemes"
	"github.com/wader/ansisvg/svgscreen"

	"flag"
	"fmt"
	"io"
	"os"
)

var fontFlag = flag.String("font", "Monaco, Lucida Console, Courier", "Font")
var fontSizeFlag = flag.Int("fontsize", 12, "Font size")
var terminalWidthFlag = flag.Int("width", 0, "Terminal width (auto)")
var characterBox = svgscreen.Dimension{Width: 7, Height: 13}
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

	a := ansidecoder.NewDecoder(os.Stdin)
	var chars []svgscreen.Char
	for {
		r, _, err := a.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		chars = append(chars, svgscreen.Char{
			Char:       string([]rune{r}),
			X:          a.X,
			Y:          a.Y,
			Foreground: a.Foreground.String(),
			Background: a.Background.String(),
			Underline:  a.Underline,
			Intensity:  a.Intensity,
		})
	}
	if *terminalWidthFlag == 0 {
		*terminalWidthFlag = a.MaxX + 1
	}

	return svgscreen.Render(
		os.Stdout,
		svgscreen.Screen{
			ColorScheme:   colorScheme.WorkbenchColorCustomizations,
			Font:          *fontFlag,
			FontSize:      *fontSizeFlag,
			CharBox:       characterBox,
			TerminalWidth: *terminalWidthFlag,
			Columns:       a.MaxX + 1,
			Lines:         a.MaxY + 1,
			Chars:         chars,
		},
	)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
