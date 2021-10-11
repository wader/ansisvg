package svgscreen

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"regexp"
	"strconv"
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
	Transparent      bool
	ForegroundColors map[string]string
	BackgroundColors map[string]string
	FontName         string
	FontSize         int
	CharacterBoxSize BoxSize
	TerminalWidth    int
	Columns          int
	Lines            int
	Chars            []Char
}

type color struct {
	R, G, B float32
}

var colorRe = regexp.MustCompile(`^#(..)(..)(..)$`)

func newColorFromHex(s string) color {
	parts := colorRe.FindStringSubmatch(s)
	if parts == nil {
		return color{}
	}
	f := func(s string) float32 { n, _ := strconv.ParseInt(s, 16, 32); return float32(n) / 255 }
	return color{
		R: f(parts[1]),
		G: f(parts[2]),
		B: f(parts[3]),
	}
}

func (c color) add(o color) color {
	clamp := func(n float32) float32 {
		if n <= 0 {
			return 0
		} else if n > 1 {
			return 1
		}
		return n
	}
	return color{
		R: clamp(c.R + o.R),
		G: clamp(c.G + o.G),
		B: clamp(c.B + o.B),
	}
}

func (c color) hex() string {
	return fmt.Sprintf("#%.2x%.2x%.2x",
		int(c.R*255),
		int(c.G*255),
		int(c.B*255),
	)
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
		"coloradd": func(a string, b string) string {
			return newColorFromHex(a).add(newColorFromHex(b)).hex()
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
