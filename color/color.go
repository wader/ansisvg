package color

import (
	"fmt"
	"regexp"
	"strconv"
)

type Color struct {
	R, G, B float32
}

var colorRe = regexp.MustCompile(`^#(..)(..)(..)$`)

func NewFromHex(s string) Color {
	parts := colorRe.FindStringSubmatch(s)
	if parts == nil {
		return Color{}
	}
	f := func(s string) float32 { n, _ := strconv.ParseInt(s, 16, 32); return float32(n) / 255 }
	return Color{
		R: f(parts[1]),
		G: f(parts[2]),
		B: f(parts[3]),
	}
}

func (c Color) Add(o Color) Color {
	clamp := func(n float32) float32 {
		if n <= 0 {
			return 0
		} else if n > 1 {
			return 1
		}
		return n
	}
	return Color{
		R: clamp(c.R + o.R),
		G: clamp(c.G + o.G),
		B: clamp(c.B + o.B),
	}
}

func (c Color) Hex() string {
	return fmt.Sprintf("#%.2x%.2x%.2x",
		int(c.R*255),
		int(c.G*255),
		int(c.B*255),
	)
}

func (c Color) ANSITriple() string {
	return fmt.Sprintf("%d:%d:%d",
		int(c.R*255),
		int(c.G*255),
		int(c.B*255),
	)
}

func (c Color) ANSIBG() string {
	return fmt.Sprintf("\x1b[48:2:%sm", c.ANSITriple())
}

func (c Color) ANSIFG() string {
	return fmt.Sprintf("\x1b[38:2:%sm", c.ANSITriple())
}
