// Package ansitosvg converts ANSI to SVG
package ansitosvg

import (
	"io"

	"github.com/wader/ansisvg/ansidecoder"
	"github.com/wader/ansisvg/colorscheme/schemes"
	"github.com/wader/ansisvg/svgscreen"
)

type BoxSize struct {
	Width  int
	Height int
}

type Options struct {
	FontName         string
	FontEmbedded     []byte
	FontRef          string
	FontSize         int
	TerminalWidth    int
	CharacterBoxSize BoxSize
	ColorScheme      string
	Transparent      bool
}

var DefaultOptions = Options{
	FontName:         "Courier",
	FontSize:         14,
	CharacterBoxSize: BoxSize{Width: 8, Height: 16},
	ColorScheme:      "Builtin Dark",
	Transparent:      false,
}

// Convert reads ANSI input from r and writes SVG to w
func Convert(r io.Reader, w io.Writer, opts Options) error {
	ad := ansidecoder.NewDecoder(r)

	lineNr := 0
	var lines []svgscreen.Line
	line := svgscreen.Line{
		Y: lineNr,
	}

	for {
		r, _, err := ad.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if r == '\n' {
			lines = append(lines, line)
			lineNr++
			line = svgscreen.Line{
				Y: lineNr,
			}
			continue
		}

		n := 1
		// normalize tab into spaces
		if r == '\t' {
			r = ' '
			n = 8 - (ad.X % 8)
		}

		for i := 0; i < n; i++ {
			line.Chars = append(line.Chars, svgscreen.Char{
				Char:       string([]rune{r}),
				X:          ad.X + i,
				Foreground: ad.Foreground.String(),
				Background: ad.Background.String(),
				Underline:  ad.Underline,
				Intensity:  ad.Intensity,
				Invert:     ad.Invert,
			})
		}
	}
	if len(line.Chars) > 0 {
		lines = append(lines, line)
	}
	terminalWidth := ad.MaxX + 1
	if opts.TerminalWidth != 0 {
		terminalWidth = opts.TerminalWidth
	}
	colorScheme, err := schemes.Load(opts.ColorScheme)
	if err != nil {
		return err
	}

	fontName := opts.FontName
	if len(opts.FontEmbedded) > 0 {
		fontName = "Embedded"
	} else if opts.FontRef != "" {
		fontName = "ExternalRef"
	}

	c := colorScheme
	return svgscreen.Render(
		w,
		svgscreen.Screen{
			Transparent:     opts.Transparent,
			ForegroundColor: c.Foreground,
			ForegroundColors: map[string]string{
				"0":  c.ANSIBlack,
				"1":  c.ANSIRed,
				"2":  c.ANSIGreen,
				"3":  c.ANSIYellow,
				"4":  c.ANSIBlue,
				"5":  c.ANSIMagenta,
				"6":  c.ANSICyan,
				"7":  c.ANSIWhite,
				"8":  c.ANSIBrightBlack,
				"9":  c.ANSIBrightRed,
				"10": c.ANSIBrightGreen,
				"11": c.ANSIBrightYellow,
				"12": c.ANSIBrightBlue,
				"13": c.ANSIBrightMagenta,
				"14": c.ANSIBrightCyan,
				"15": c.ANSIBrightWhite,
			},
			BackgroundColor: c.Background,
			BackgroundColors: map[string]string{
				"0":  c.ANSIBlack,
				"1":  c.ANSIRed,
				"2":  c.ANSIGreen,
				"3":  c.ANSIYellow,
				"4":  c.ANSIBlue,
				"5":  c.ANSIMagenta,
				"6":  c.ANSICyan,
				"7":  c.ANSIWhite,
				"8":  c.ANSIBrightBlack,
				"9":  c.ANSIBrightRed,
				"10": c.ANSIBrightYellow,
				"11": c.ANSIBrightYellow,
				"12": c.ANSIBrightBlue,
				"13": c.ANSIBrightMagenta,
				"14": c.ANSIBrightCyan,
				"15": c.ANSIBrightWhite,
			},
			FontName:     fontName,
			FontEmbedded: opts.FontEmbedded,
			FontRef:      opts.FontRef,
			FontSize:     opts.FontSize,
			CharacterBoxSize: svgscreen.BoxSize{
				Width:  opts.CharacterBoxSize.Width,
				Height: opts.CharacterBoxSize.Height,
			},
			TerminalWidth: terminalWidth,
			Columns:       ad.MaxX + 1,
			NrLines:       ad.MaxY + 1,
			Lines:         lines,
		},
	)
}
