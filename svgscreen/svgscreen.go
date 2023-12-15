// Package svgscreen implements a fixed font terminal screen using SVG
package svgscreen

import (
	_ "embed"
	"encoding/base64"
	"html/template"
	"io"
	"strings"

	"github.com/wader/ansisvg/color"
)

//go:embed template.svg
var templateSVG string

type Char struct {
	Char       string
	X          int
	Foreground string
	Background string
	Underline  bool
	Intensity  bool
	Invert     bool
}

type Line struct {
	Y     int
	Chars []Char
}

type BoxSize struct {
	Width  int
	Height int
}

type Screen struct {
	Transparent      bool
	ForegroundColor  string
	ForegroundColors map[string]string
	BackgroundColor  string
	BackgroundColors map[string]string
	FontName         string
	FontEmbedded     []byte
	FontRef          string
	FontSize         int
	CharacterBoxSize BoxSize
	TerminalWidth    int
	Columns          int
	NrLines          int
	Lines            []Line
}

func Render(w io.Writer, s Screen) error {
	t := template.New("")
	t.Funcs(template.FuncMap{
		"add":          func(a int, b int) int { return a + b },
		"mul":          func(a int, b int) int { return a * b },
		"hasprefix":    strings.HasPrefix,
		"iswhitespace": func(a string) bool { return strings.TrimSpace(a) == "" },
		"coloradd": func(a string, b string) string {
			return color.NewFromHex(a).Add(color.NewFromHex(b)).Hex()
		},
		"base64": func(bs []byte) string { return base64.RawStdEncoding.EncodeToString(bs) },
	})

	// remove unused background colors
	backgroundColorsUsed := map[string]string{}
	for _, l := range s.Lines {
		for i, c := range l.Chars {
			if c.Invert {
				c.Background, c.Foreground = c.Foreground, c.Background
				if c.Background == "" {
					c.Background = s.ForegroundColor
				}
				if c.Foreground == "" {
					c.Foreground = s.BackgroundColor
				}
				l.Chars[i] = c
			}

			if c.Background == "" {
				continue
			}
			if strings.HasPrefix(c.Background, "#") {
				backgroundColorsUsed[c.Background] = c.Background
			} else {
				backgroundColorsUsed[c.Background] = s.BackgroundColors[c.Background]
			}
		}
	}
	s.BackgroundColors = backgroundColorsUsed

	t, err := t.Parse(templateSVG)
	if err != nil {
		return err
	}
	if err = t.ExecuteTemplate(w, "", s); err != nil {
		return err
	}

	return nil
}
