// Package colorscheme has defintions for a VSCode color scheme
package colorscheme

import (
	"fmt"
	"strings"

	"github.com/wader/ansisvg/color"
)

type WorkbenchColorCustomizations struct {
	Foreground          string `json:"terminal.foreground"`
	Background          string `json:"terminal.background"`
	ANSIBlack           string `json:"terminal.ansiBlack"`
	ANSIBlue            string `json:"terminal.ansiBlue"`
	ANSICyan            string `json:"terminal.ansiCyan"`
	ANSIGreen           string `json:"terminal.ansiGreen"`
	ANSIMagenta         string `json:"terminal.ansiMagenta"`
	ANSIRed             string `json:"terminal.ansiRed"`
	ANSIWhite           string `json:"terminal.ansiWhite"`
	ANSIYellow          string `json:"terminal.ansiYellow"`
	ANSIBrightBlack     string `json:"terminal.ansiBrightBlack"`
	ANSIBrightBlue      string `json:"terminal.ansiBrightBlue"`
	ANSIBrightCyan      string `json:"terminal.ansiBrightCyan"`
	ANSIBrightGreen     string `json:"terminal.ansiBrightGreen"`
	ANSIBrightMagenta   string `json:"terminal.ansiBrightMagenta"`
	ANSIBrightRed       string `json:"terminal.ansiBrightRed"`
	ANSIBrightWhite     string `json:"terminal.ansiBrightWhite"`
	ANSIBrightYellow    string `json:"terminal.ansiBrightYellow"`
	SelectionBackground string `json:"terminal.selectionBackground"`
	CursorForeground    string `json:"terminalCursor.foreground"`
}

type VSCodeColorScheme struct {
	WorkbenchColorCustomizations WorkbenchColorCustomizations `json:"workbench.colorCustomizations"`
}

func (w WorkbenchColorCustomizations) ANSIDemo(s string) string {
	b := color.NewFromHex(w.Background)
	f := color.NewFromHex(w.Foreground)
	var sb strings.Builder
	for _, c := range []string{
		w.ANSIBlack,
		w.ANSIBlue,
		w.ANSICyan,
		w.ANSIGreen,
		w.ANSIMagenta,
		w.ANSIRed,
		w.ANSIWhite,
		w.ANSIYellow,
		w.ANSIBrightBlack,
		w.ANSIBrightBlue,
		w.ANSIBrightCyan,
		w.ANSIBrightGreen,
		w.ANSIBrightMagenta,
		w.ANSIBrightRed,
		w.ANSIBrightWhite,
		w.ANSIBrightYellow,
		w.SelectionBackground,
		w.CursorForeground,
	} {
		sb.WriteString(fmt.Sprintf("%s  \x1b[0m", color.NewFromHex(c).ANSIBG()))
	}

	return fmt.Sprintf("%s%s%s%s\x1b[0m",
		b.ANSIBG(),
		f.ANSIFG(),
		s,
		sb.String(),
	)
}
