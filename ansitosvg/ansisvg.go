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
	FontSize         int
	TerminalWidth    int
	CharacterBoxSize BoxSize
	ColorScheme      string
	Transparent      bool
	CompactByStyle   bool
}

var DefaultOptions = Options{
	FontName:         "Courier",
	FontSize:         14,
	CharacterBoxSize: BoxSize{Width: 8, Height: 16},
	ColorScheme:      "Builtin Dark",
	Transparent:      false,
	CompactByStyle:   false,
}

// Convert reads ANSI input from r and writes SVG to w
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

		if r == '\n' {
			continue
		}

		n := 1
		// normalize tab into spaces
		if r == '\t' {
			r = ' '
			n = 8 - (ad.X % 8)
		}
		for i := 0; i < n; i++ {
			chars = append(chars, svgscreen.Char{
				Char:       string([]rune{r}),
				X:          ad.X + i,
				Y:          ad.Y,
				Foreground: ad.Foreground.String(),
				Background: ad.Background.String(),
				Underline:  ad.Underline,
				Intensity:  ad.Intensity,
			})
		}
	}

	if opts.CompactByStyle {

		// log.Println("chars:")
		// for _, c := range chars {
		// 	log.Printf("c: %#+v\n", c)
		// }

		// TODO: check len
		var compactedChars []svgscreen.Char
		var p svgscreen.Char
		hasP := false
		for i, c := range chars {
			if i == 0 {
				p = c
				hasP = true
			} else if p.Y == c.Y && p.HasSameStyle(c) {
				p.Char += c.Char
				hasP = true
			} else {
				if hasP {
					compactedChars = append(compactedChars, p)
				}
				p = c
			}
		}
		compactedChars = append(compactedChars, p)

		chars = compactedChars

		// log.Println("compactedChars:")
		// for _, c := range compactedChars {
		// 	log.Printf("c: %#+v\n", c)
		// }

	}

	terminalWidth := ad.MaxX + 1
	if opts.TerminalWidth != 0 {
		terminalWidth = opts.TerminalWidth
	}
	colorScheme, err := schemes.Load(opts.ColorScheme)
	if err != nil {
		return err
	}

	c := colorScheme
	return svgscreen.Render(
		w,
		svgscreen.Screen{
			Transparent:     opts.Transparent,
			ForegroundColor: c.Foreground,
			ForegroundColors: map[string]string{
				"30": c.ANSIBlack,
				"31": c.ANSIRed,
				"32": c.ANSIGreen,
				"33": c.ANSIYellow,
				"34": c.ANSIBlue,
				"35": c.ANSIMagenta,
				"36": c.ANSICyan,
				"37": c.ANSIWhite,
				"90": c.ANSIBrightBlack,
				"91": c.ANSIBrightRed,
				"92": c.ANSIBrightGreen,
				"93": c.ANSIBrightYellow,
				"94": c.ANSIBrightBlue,
				"95": c.ANSIBrightMagenta,
				"96": c.ANSIBrightCyan,
				"97": c.ANSIBrightWhite,
			},
			BackgroundColor: c.Background,
			BackgroundColors: map[string]string{
				"40":  c.ANSIBlack,
				"41":  c.ANSIRed,
				"42":  c.ANSIGreen,
				"43":  c.ANSIYellow,
				"44":  c.ANSIBlue,
				"45":  c.ANSIMagenta,
				"46":  c.ANSICyan,
				"47":  c.ANSIWhite,
				"100": c.ANSIBrightBlack,
				"101": c.ANSIBrightRed,
				"102": c.ANSIBrightYellow,
				"103": c.ANSIBrightYellow,
				"104": c.ANSIBrightBlue,
				"105": c.ANSIBrightMagenta,
				"106": c.ANSIBrightCyan,
				"107": c.ANSIBrightWhite,
			},
			FontName: opts.FontName,
			FontSize: opts.FontSize,
			CharacterBoxSize: svgscreen.BoxSize{
				Width:  opts.CharacterBoxSize.Width,
				Height: opts.CharacterBoxSize.Height,
			},
			TerminalWidth: terminalWidth,
			Columns:       ad.MaxX + 1,
			Lines:         ad.MaxY + 1,
			Chars:         chars,
		},
	)
}
