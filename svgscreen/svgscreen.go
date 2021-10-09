package svgscreen

import (
	_ "embed"
	"html/template"
	"io"
	"strings"
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

type BoxSize struct {
	Width  int
	Height int
}

type Screen struct {
	ForegroundColors [256]string
	BackgroundColors [256]string
	Font             string
	FontSize         int
	CharacterBoxSize BoxSize
	TerminalWidth    int
	Columns          int
	Lines            int
	Chars            []Char
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
