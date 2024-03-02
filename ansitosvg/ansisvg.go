// Package ansitosvg converts ANSI to SVG
package ansitosvg

import (
	"io"

	"github.com/wader/ansisvg/ansidecoder"
	"github.com/wader/ansisvg/colorscheme/schemes"
	"github.com/wader/ansisvg/svgscreen"
	"github.com/wader/ansisvg/svgscreen/xydim"
)

type Options struct {
	FontName      string
	FontEmbedded  []byte
	FontRef       string
	FontSize      int
	TerminalWidth int
	CharBoxSize   xydim.XyDimInt
	ColorScheme   string
	Transparent   bool
	GridMode      bool
}

var DefaultOptions = Options{
	FontName:    "Courier",
	FontSize:    14,
	CharBoxSize: xydim.XyDimInt{X: 0, Y: 0},
	ColorScheme: "Builtin Dark",
	Transparent: false,
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
				Char:          string([]rune{r}),
				X:             ad.X + i,
				Foreground:    ad.Foreground.String(),
				Background:    ad.Background.String(),
				Underline:     ad.Underline,
				Intensity:     ad.Intensity,
				Invert:        ad.Invert,
				Italic:        ad.Italic,
				Strikethrough: ad.Strikethrough,
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
	s := svgscreen.Screen{
		Transparent: opts.Transparent,
		Foreground: svgscreen.ColorMap{
			Default: c.Foreground,
		},
		Background: svgscreen.ColorMap{
			Default: c.Background,
		},
		ANSIColors: [16]string{
			c.ANSIBlack,
			c.ANSIRed,
			c.ANSIGreen,
			c.ANSIYellow,
			c.ANSIBlue,
			c.ANSIMagenta,
			c.ANSICyan,
			c.ANSIWhite,
			c.ANSIBrightBlack,
			c.ANSIBrightRed,
			c.ANSIBrightGreen,
			c.ANSIBrightYellow,
			c.ANSIBrightBlue,
			c.ANSIBrightMagenta,
			c.ANSIBrightCyan,
			c.ANSIBrightWhite,
		},
		Dom: svgscreen.SvgDom{
			FontName:     fontName,
			FontEmbedded: opts.FontEmbedded,
			FontRef:      opts.FontRef,
			FontSize:     opts.FontSize,
		},
		CharacterBoxSize: opts.CharBoxSize,
		TerminalWidth:    terminalWidth,
		Columns:          ad.MaxX + 1,
		NrLines:          ad.MaxY + 1,
		Lines:            lines,
		GridMode:         opts.GridMode,
	}
	return s.Render(w)
}
