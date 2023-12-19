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
	Char          string
	X             int
	Foreground    string
	Background    string
	Underline     bool
	Intensity     bool
	Invert        bool
	Italic        bool
	Strikethrough bool
}

type Line struct {
	Y     int
	Chars []Char
}

type BoxSize struct {
	Width  int
	Height int
}

type TextSpan struct {
	Style      template.CSS
	Decoration template.CSS
	Content    string
}

type TextElement struct {
	Y         int
	TextSpans []TextSpan
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
	TextElements     []TextElement
}

// Resolve color from string (either # prefixed hex value or index into lookup table)
func ResolveColor(c string, lookup map[string]string) string {
	if strings.HasPrefix(c, "#") {
		return c
	}
	return lookup[c]
}

// Convert a line into a <text> element
// fc gives (foregroundcolor, content) of a char
func LineToTextElement(s Screen, l Line, fc func(Char) TextSpan) TextElement {
	var t []TextSpan
	currentSpan := TextSpan{
		Style:      "",
		Decoration: "",
		Content:    "",
	}

	appendSpan := func() {
		if currentSpan.Content == "" {
			return
		}
		t = append(t, currentSpan)
	}
	for _, c := range l.Chars {
		tempSpan := fc(c)
		if tempSpan.Style != currentSpan.Style || tempSpan.Decoration != currentSpan.Decoration {
			appendSpan()
			currentSpan = tempSpan
			continue
		}
		currentSpan.Content += tempSpan.Content
	}
	appendSpan()

	// remove trailing whitespace
	for len(t) > 0 && strings.TrimSpace(t[len(t)-1].Content) == "" {
		t = t[:len(t)-1]
	}

	return TextElement{
		Y:         l.Y * s.CharacterBoxSize.Height,
		TextSpans: t,
	}
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
		}
	}

	// Render the whole background first. Then it will not be selected when copy/pasting from rendered SVG
	for _, l := range s.Lines {
		bg := LineToTextElement(s, l, func(c Char) TextSpan {
			if c.Background == "" {
				return TextSpan{
					Style:   "",
					Content: " ",
				}
			} else {
				return TextSpan{
					Style:   template.CSS("fill: " + ResolveColor(c.Background, s.BackgroundColors)),
					Content: "█",
				}
			}
		})
		if len(bg.TextSpans) > 0 {
			s.TextElements = append(s.TextElements, bg)
		}
	}

	// Then render the foreground
	for _, l := range s.Lines {
		fg := LineToTextElement(s, l, func(c Char) TextSpan {
			var styles []string
			deco := ""

			if c.Foreground != "" {
				styles = append(styles, "fill:"+ResolveColor(c.Foreground, s.ForegroundColors))
			}
			if c.Intensity {
				styles = append(styles, "font-weight:bold")
			}
			if c.Italic {
				styles = append(styles, "font-style:italic")
			}
			if c.Underline {
				deco = "underline"
			} else if c.Strikethrough {
				deco = "line-through"
			}

			return TextSpan{
				Style:      template.CSS(strings.Join(styles, "; ")),
				Decoration: template.CSS(deco),
				Content:    c.Char,
			}
		})
		if len(fg.TextSpans) > 0 {
			s.TextElements = append(s.TextElements, fg)
		}
	}

	t, err := t.Parse(templateSVG)
	if err != nil {
		return err
	}
	if err = t.ExecuteTemplate(w, "", s); err != nil {
		return err
	}

	return nil
}
