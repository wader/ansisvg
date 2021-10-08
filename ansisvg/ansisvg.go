package ansisvg

import (
	"io"

	"github.com/wader/ansisvg/ansidecoder"
	"github.com/wader/ansisvg/colorscheme"
	"github.com/wader/ansisvg/svgscreen"
)

type Dimension struct {
	Width  int
	Height int
}

type Options struct {
	Font          string
	FontSize      int
	TerminalWidth int
	CharacterBox  Dimension
	ColorScheme   colorscheme.WorkbenchColorCustomizations
}

func Convert(r io.Reader, w io.Writer, opts Options) error {
	ad := ansidecoder.NewDecoder(r)
	var chars []svgscreen.Char
	for {
		r, _, err := ad.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		chars = append(chars, svgscreen.Char{
			Char:       string([]rune{r}),
			X:          ad.X,
			Y:          ad.Y,
			Foreground: ad.Foreground.String(),
			Background: ad.Background.String(),
			Underline:  ad.Underline,
			Intensity:  ad.Intensity,
		})
	}
	terminalWidth := ad.MaxX + 1
	if opts.TerminalWidth != 0 {
		terminalWidth = opts.TerminalWidth
	}

	c := opts.ColorScheme
	return svgscreen.Render(
		w,
		svgscreen.Screen{
			ForegroundColors: [256]string{
				0:  c.Foreground,
				30: c.ANSIBlack,
				31: c.ANSIRed,
				32: c.ANSIGreen,
				33: c.ANSIYellow,
				34: c.ANSIBlue,
				35: c.ANSIMagenta,
				36: c.ANSICyan,
				37: c.ANSIWhite,
				90: c.ANSIBrightBlack,
				91: c.ANSIBrightRed,
				92: c.ANSIBrightGreen,
				93: c.ANSIBrightYellow,
				94: c.ANSIBrightBlue,
				95: c.ANSIBrightMagenta,
				96: c.ANSIBrightCyan,
				97: c.ANSIBrightWhite,
			},
			BackgroundColors: [256]string{
				0:   c.Background,
				40:  c.ANSIBlack,
				41:  c.ANSIRed,
				42:  c.ANSIYellow,
				43:  c.ANSIYellow,
				44:  c.ANSIBlue,
				45:  c.ANSIMagenta,
				46:  c.ANSICyan,
				47:  c.ANSIWhite,
				100: c.ANSIBrightBlack,
				101: c.ANSIBrightRed,
				102: c.ANSIBrightYellow,
				103: c.ANSIBrightYellow,
				104: c.ANSIBrightBlue,
				105: c.ANSIBrightMagenta,
				106: c.ANSIBrightCyan,
				107: c.ANSIBrightWhite,
			},
			Font:     opts.Font,
			FontSize: opts.FontSize,
			CharacterBox: svgscreen.Dimension{
				Width:  opts.CharacterBox.Width,
				Height: opts.CharacterBox.Height,
			},
			TerminalWidth: terminalWidth,
			Columns:       ad.MaxX + 1,
			Lines:         ad.MaxY + 1,
			Chars:         chars,
		},
	)
}
