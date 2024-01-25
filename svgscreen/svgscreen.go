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
	X       string
	Class   string
	Content string
}

type textElement struct {
	Y         string
	TextSpans []textSpan
}

type bgRect struct {
	X      string
	Y      string
	Width  string
	Height string
	Color  string
}

type SvgDom struct {
	Width          string
	Height         string
	FontName       string
	FontEmbedded   []byte
	FontRef        string
	FontSize       int
	FgCustomColors []string
	BgCustomColors []string
	BgRects        []bgRect
	TextElements   []textElement
	ClassesUsed    struct {
		Bold          bool
		Italic        bool
		Underline     bool
		Strikethrough bool
	}
}

type ColorMap struct {
	Default   string
	Custom    map[string]int
	ANSIUsed  [16]bool
	DomPrefix string
}

type Screen struct {
	Transparent bool
	Foreground  ColorMap
	Background  ColorMap
	ANSIColors  [16]string

	CharacterBoxSize BoxSize
	TerminalWidth    int
	Columns          int
	NrLines          int
	Lines            []Line
	GridMode         bool
	Dom              SvgDom
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
func (s *Screen) resolveColor(c string, cmap *ColorMap) string {
	if c == "" || c == cmap.Default {
		return ""
	}

	if !strings.HasPrefix(c, "#") {
		// standard ANSI color
		idx, _ := strconv.Atoi(c)
		cmap.ANSIUsed[idx] = true
		return cmap.DomPrefix + "a" + c
	}
	// custom color. update lookup table if necessary
	colIdx, present := cmap.Custom[c]
	if present {
		return cmap.DomPrefix + "c" + strconv.Itoa(colIdx)
	}
	colIdx = len(cmap.Custom)
	cmap.Custom[c] = colIdx
	return cmap.DomPrefix + "c" + strconv.Itoa(colIdx)
}

func (s *Screen) charToFgText(c Char) textSpan {
	var classes []string

	if c.Intensity {
		classes = append(classes, "bold")
		s.Dom.ClassesUsed.Bold = true
	}
	if c.Italic {
		classes = append(classes, "italic")
		s.Dom.ClassesUsed.Italic = true
	}
	if c.Underline {
		classes = append(classes, "underline")
		s.Dom.ClassesUsed.Underline = true
	} else if c.Strikethrough {
		classes = append(classes, "strikethrough")
		s.Dom.ClassesUsed.Strikethrough = true
	}

	color := s.resolveColor(c.Foreground, &s.Foreground)
	if color != "" {
		classes = append(classes, color)
	}

	return textSpan{
		Class:   strings.Join(classes, " "),
		Content: c.Char,
	}
}

// Convert a line into a textElement
// fc gives textSpan for a single char
func (s *Screen) lineToTextElement(l Line) textElement {
	var t []textSpan
	currentSpan := textSpan{
		Class:   "",
		Content: "",
	}

	appendSpan := func() {
		if currentSpan.Content == "" {
			return
		}
		t = append(t, currentSpan)
	}
	for col, c := range l.Chars {
		newSpan := s.charToFgText(c)
		if s.GridMode {
			// In grid mode, set X coordinate for each text span
			newSpan.X = s.columnCoordinate(col)
			// In grid mode, we never consolidate
			appendSpan()
			currentSpan = newSpan
			continue
		}
		// Don't consolidate if class is changing, but ignore whitespace
		if newSpan.Class != currentSpan.Class && strings.TrimSpace(newSpan.Content) != "" {
			appendSpan()
			currentSpan = newSpan
			continue
		}
		// Consolidate new content with previous one.
		currentSpan.Content += newSpan.Content
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
					c.Background = s.Foreground.Default
				}
				if c.Foreground == "" {
					c.Foreground = s.Background.Default
				}
				l.Chars[i] = c
				c.Invert = false
			}
		}
	}
}

func (s *Screen) setupBgRects() {
	// Set up background rects
	for y, l := range s.Lines {
		type tmpRect struct {
			x     int
			w     int
			color string
		}
		currentRect := tmpRect{x: 0, w: 0, color: ""}

		appendRect := func() {
			if currentRect.color == "" {
				return
			}
			s.Dom.BgRects = append(s.Dom.BgRects, bgRect{
				X:      s.columnCoordinate(currentRect.x),
				Y:      s.rowCoordinate(float32(y)),
				Width:  s.columnCoordinate(currentRect.w),
				Height: s.rowCoordinate(1),
				Color:  currentRect.color,
			})
		}
		for x, c := range l.Chars {
			if c.Background == "" || c.Background == s.Background.Default {
				continue
			}
			newRect := tmpRect{x: x, w: 1, color: s.resolveColor(c.Background, &s.Background)}

			if newRect.x != (currentRect.x+currentRect.w) || newRect.color != currentRect.color {
				appendRect()
				currentRect = newRect
				continue
			}

			currentRect.w += 1
		}
		appendRect()
	}
}

func setupCustomColors(revLookup map[string]int, clsTable *[]string) {
	result := make([]string, len(revLookup))
	for k, v := range revLookup {
		result[v] = k
	}
	*clsTable = result
}

func (s *Screen) Render(w io.Writer) error {
	t := template.New("")
	t.Funcs(template.FuncMap{
		"base64": func(bs []byte) string { return base64.RawStdEncoding.EncodeToString(bs) },
		"anyColorUsed": func(arr [16]bool) bool {
			for _, value := range arr {
				if value {
					return true
				}
			}
			return false
		},
	})

	s.Foreground.DomPrefix = "f"
	s.Background.DomPrefix = "b"
	s.Foreground.Custom = map[string]int{}
	s.Background.Custom = map[string]int{}

	// Set SVG size
	s.Dom.Width = s.columnCoordinate(s.TerminalWidth)
	s.Dom.Height = s.rowCoordinate(float32(s.NrLines))

	s.handleColorInversion()
	s.setupBgRects()

	// Set up text elements
	for _, l := range s.Lines {
		fg := s.lineToTextElement(l)
		if len(fg.TextSpans) > 0 {
			s.Dom.TextElements = append(s.Dom.TextElements, fg)
		}
	}

	setupCustomColors(s.Foreground.Custom, &s.Dom.FgCustomColors)
	setupCustomColors(s.Background.Custom, &s.Dom.BgCustomColors)

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
