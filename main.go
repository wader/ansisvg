package main

import (
	"github.com/wader/ansisvg/ansidecoder"

	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"strconv"
	"strings"
)

//go:embed template.svg
var svgTemplate string

//go:embed colorschemas/*.json
var colorSchemas embed.FS

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

type Char struct {
	Char       string
	X          int
	Y          int
	Foreground string
	Background string
	Underline  bool
	Intensity  bool
}

type workbenchColorCustomizations struct {
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

type vsCodeColorScheme struct {
	WorkbenchColorCustomizations workbenchColorCustomizations `json:"workbench.colorCustomizations"`
}

func run() error {
	flag.Parse()

	a := ansidecoder.NewDecoder(os.Stdin)
	var chars []Char
	for {
		r, _, err := a.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		chars = append(chars, Char{
			Char:       string([]rune{r}),
			X:          a.X,
			Y:          a.Y,
			Foreground: a.Foreground.String(),
			Background: a.Background.String(),
			Underline:  a.Underline,
			Intensity:  a.Intensity,
		})
	}

	var colorScheme vsCodeColorScheme

	f, err := colorSchemas.Open("colorschemas/" + *colorSchemeFlag + ".json")
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&colorScheme); err != nil {
		return err
	}

	if *terminalWidthFlag == 0 {
		*terminalWidthFlag = a.MaxX + 1
	}

	t := template.New("")
	t.Funcs(template.FuncMap{
		"add": func(a int, b int) int {
			return a + b
		},
		"mul": func(a int, b int) int {
			return a * b
		},
		"hasprefix": strings.HasPrefix,
		"iswhitespace": func(a string) bool {
			return strings.TrimSpace(a) == ""
		},
	})

	t, err = t.Parse(svgTemplate)
	if err != nil {
		return err
	}
	if err = t.ExecuteTemplate(os.Stdout, "", struct {
		ColorScheme   workbenchColorCustomizations
		Font          string
		FontSize      int
		CharBox       dimension
		TerminalWidth int
		Columns       int
		Lines         int
		Chars         []Char
	}{
		ColorScheme:   colorScheme.WorkbenchColorCustomizations,
		Font:          *fontFlag,
		FontSize:      *fontSizeFlag,
		CharBox:       characterBox,
		TerminalWidth: *terminalWidthFlag,
		Columns:       a.MaxX + 1,
		Lines:         a.MaxY + 1,
		Chars:         chars,
	}); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
