package svgscreen

import (
	_ "embed"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"

	"github.com/wader/ansisvg/colorscheme"
)

//go:embed template.svg
var templateSVG string

type Char struct {
	Char       string
	X          int
	Y          int
	Foreground string
	Background string
	Underline  bool
	Intensity  bool
}

type Dimension struct {
	Width  int
	Height int
}

func (d *Dimension) String() string {
	return fmt.Sprintf("%dx%d", d.Width, d.Height)
}

func (d *Dimension) Set(s string) error {
	parts := strings.Split(s, "x")
	if len(parts) != 2 {
		return fmt.Errorf("must be WxH")
	}
	d.Width, _ = strconv.Atoi(parts[0])
	d.Height, _ = strconv.Atoi(parts[1])
	return nil
}

type Screen struct {
	ColorScheme   colorscheme.WorkbenchColorCustomizations
	Font          string
	FontSize      int
	CharBox       Dimension
	TerminalWidth int
	Columns       int
	Lines         int
	Chars         []Char
}

func Render(w io.Writer, s Screen) error {
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

	t, err := t.Parse(templateSVG)
	if err != nil {
		return err
	}
	if err = t.ExecuteTemplate(w, "", s); err != nil {
		return err
	}

	return nil
}
