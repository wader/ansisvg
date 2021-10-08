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
			ForegroundColors: [256]string{
				0:  colorScheme.Foreground,
				30: colorScheme.ANSIBlack,
				31: colorScheme.ANSIRed,
				32: colorScheme.ANSIGreen,
				33: colorScheme.ANSIYellow,
				34: colorScheme.ANSIBlue,
				35: colorScheme.ANSIMagenta,
				36: colorScheme.ANSICyan,
				37: colorScheme.ANSIWhite,
				90: colorScheme.ANSIBrightBlack,
				91: colorScheme.ANSIBrightRed,
				92: colorScheme.ANSIBrightGreen,
				93: colorScheme.ANSIBrightYellow,
				94: colorScheme.ANSIBrightBlue,
				95: colorScheme.ANSIBrightMagenta,
				96: colorScheme.ANSIBrightCyan,
				97: colorScheme.ANSIBrightWhite,
			},
			BackgroundColors: [256]string{
				0:   colorScheme.Background,
				40:  colorScheme.ANSIBlack,
				41:  colorScheme.ANSIRed,
				42:  colorScheme.ANSIYellow,
				43:  colorScheme.ANSIYellow,
				44:  colorScheme.ANSIBlue,
				45:  colorScheme.ANSIMagenta,
				46:  colorScheme.ANSICyan,
				47:  colorScheme.ANSIWhite,
				100: colorScheme.ANSIBrightBlack,
				101: colorScheme.ANSIBrightRed,
				102: colorScheme.ANSIBrightYellow,
				103: colorScheme.ANSIBrightYellow,
				104: colorScheme.ANSIBrightBlue,
				105: colorScheme.ANSIBrightMagenta,
				106: colorScheme.ANSIBrightCyan,
				107: colorScheme.ANSIBrightWhite,
			},
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
