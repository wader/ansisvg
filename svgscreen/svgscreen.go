// Package svgscreen implements a fixed font terminal screen using SVG
package svgscreen

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"strconv"
	"strings"
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

type textSpan struct {
	X          string
	Style      template.CSS
	Decoration template.CSS
	Content    string
}

type textElement struct {
	Y         string
	TextSpans []textSpan
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
	TextElements     []textElement
	GridMode         bool
	SvgWidth         string
	SvgHeight        string
}

func (s *Screen) columnCoordinate(col int) string {
	unit := "ch"
	if s.CharacterBoxSize.Width > 0 {
		unit = "px"
		col *= s.CharacterBoxSize.Width
	}
	return strconv.Itoa(col) + unit
}

func (s *Screen) rowCoordinate(row float32) string {
	unit := "em"
	if s.CharacterBoxSize.Height > 0 {
		unit = "px"
		row *= float32(s.CharacterBoxSize.Height)
	}
	return fmt.Sprintf("%g%s", row, unit)
}

// Resolve color from string (either # prefixed hex value or index into lookup table)
func resolveColor(c string, lookup map[string]string) string {
	if strings.HasPrefix(c, "#") {
		return c
	}
	return lookup[c]
}

// Convert a line into a textElement
// fc gives textSpan for a single char
func (s *Screen) lineToTextElement(l Line, fc func(Char) textSpan) textElement {
	var t []textSpan
	currentSpan := textSpan{
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
	for col, c := range l.Chars {
		tempSpan := fc(c)
		if s.GridMode {
			// If in grid mode, set X coordinate for each text span
			tempSpan.X = s.columnCoordinate(col)
		}
		// Consolidate tempSpan to currentSpan only if not in grid mode and both spans have same style
		if s.GridMode || tempSpan.Style != currentSpan.Style || tempSpan.Decoration != currentSpan.Decoration {
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

	return textElement{
		Y:         s.rowCoordinate(float32(l.Y) + 0.5),
		TextSpans: t,
	}
}

func (s *Screen) handleColorInversion() {
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
				c.Invert = false
			}
		}
	}
}

func (s *Screen) charToBgText(c Char) textSpan {
	if c.Background == "" {
		return textSpan{
			Style:   "",
			Content: " ",
		}
	} else {
		return textSpan{
			Style:   template.CSS("fill:" + resolveColor(c.Background, s.BackgroundColors)),
			Content: "█",
		}
	}
}

func (s *Screen) charToFgText(c Char) textSpan {
	var styles []string
	deco := ""

	if c.Foreground != "" {
		styles = append(styles, "fill:"+resolveColor(c.Foreground, s.ForegroundColors))
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

	return textSpan{
		Style:      template.CSS(strings.Join(styles, ";")),
		Decoration: template.CSS(deco),
		Content:    c.Char,
	}
}

func (s *Screen) Render(w io.Writer) error {
	t := template.New("")
	t.Funcs(template.FuncMap{
		"base64": func(bs []byte) string { return base64.RawStdEncoding.EncodeToString(bs) },
	})

	// Set SVG size
	s.SvgWidth = s.columnCoordinate(s.TerminalWidth)
	s.SvgHeight = s.rowCoordinate(float32(s.NrLines))

	s.handleColorInversion()

	// Render the whole background first. Then it will not be selected when copy/pasting from rendered SVG
	for _, l := range s.Lines {
		bg := s.lineToTextElement(l, s.charToBgText)
		if len(bg.TextSpans) > 0 {
			s.TextElements = append(s.TextElements, bg)
		}
	}

	// Then render the foreground
	for _, l := range s.Lines {
		fg := s.lineToTextElement(l, s.charToFgText)
		if len(fg.TextSpans) > 0 {
			s.TextElements = append(s.TextElements, fg)
		}
	}

	// Create SVG from template
	t, err := t.Parse(templateSVG)
	if err != nil {
		return err
	}
	if err = t.ExecuteTemplate(w, "", s); err != nil {
		return err
	}

	return nil
}
