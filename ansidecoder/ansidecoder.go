// Package ansidecoder implements a ANSI decoder that returns runes and
// keeps track of cursor position and styling.
package ansidecoder

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

type State int

const (
	StateCopy       State = iota
	StateSeenESC          // Seen ESC
	StateCSI              // Control Sequence Inducer ESC [
	StateOSC              // Operating System Command ESC ]
	StateOSCSeenESC       // Operating System Command ESC ] ... ESC
)

type codeRanges [2]int

func (cr codeRanges) Is(c int) bool {
	if c >= cr[0] && c <= cr[1] {
		return true
	}
	return false
}

var sgrReset = codeRanges{0, 0}
var sgrIncreaseIntensity = codeRanges{1, 1}
var sgrNormal = codeRanges{22, 22}
var sgrForeground = codeRanges{30, 37}
var sgrForegroundBright = codeRanges{90, 97}
var sgrForegroundRGB = codeRanges{38, 38}
var sgrForegroundDefault = codeRanges{39, 39}
var sgrBackground = codeRanges{40, 47}
var sgrBackgroundBright = codeRanges{100, 107}
var sgrBackgroundRGB = codeRanges{48, 48}
var sgrBackgroundDefault = codeRanges{49, 49}
var sgrUnderlineOn = codeRanges{4, 4}
var sgrUnderlineOff = codeRanges{24, 24}
var sgrInvertOn = codeRanges{7, 7}
var sgrInvertOff = codeRanges{27, 27}

const ESCRune = rune('\x1b')
const BELRune = rune('\x07')
const SGRByte = 'm' // Select Graphic Rendition
const FinalBytes = "@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~)"

type Color struct {
	N   int
	RGB []int
}

func (c Color) String() string {
	if len(c.RGB) != 0 {
		return fmt.Sprintf("#%.2x%.2x%.2x", c.RGB[0], c.RGB[1], c.RGB[2])
	}
	if c.N != -1 {
		return fmt.Sprintf("%d", c.N)
	}
	return ""
}

type Decoder struct {
	// state of last returned rune
	X          int
	Y          int
	Foreground Color
	Background Color
	Underline  bool
	Intensity  bool
	Invert     bool

	MaxX  int
	MaxY  int
	State State

	// next coordinate
	nx        int
	ny        int
	readBuf   *bufio.Reader
	paramsBuf *bytes.Buffer
}

// NewDecoder returns new ANSI decoder that is a io.RuneReader. See ReadRune for details.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		Foreground: Color{N: -1},
		Background: Color{N: -1},
		readBuf:    bufio.NewReader(r),
		paramsBuf:  &bytes.Buffer{},
	}
}

func intsToColor(fo int, bo int, cs []int) (Color, int) {
	if len(cs) == 0 {
		return Color{N: -1}, 0
	}
	switch {
	case cs[0] == 2 && len(cs) >= 4: // 2;r;g;b
		return Color{RGB: append([]int{}, cs[1:4]...)}, 4
	case cs[0] == 5 && len(cs) >= 2: // 5;n
		n := cs[1]
		switch {
		case n >= 0 && n <= 15:
			// 0-  7:  standard colors (as in ESC [ 30–37 m)
			// 8- 15:  high intensity colors (as in ESC [ 90–97 m)
			return Color{N: n}, 2
		case n >= 16 && n <= 231:
			// 16-231:  6 × 6 × 6 cube (216 colors): 16 + 36 × r + 6 × g + b (0 ≤ r, g, b ≤ 5)
			// TODO: not tested
			n -= 16
			r := n / 36
			n %= 36
			g := n / 6
			n %= 6
			b := n

			// iterm2 mapping of 0-5 -> 0-255 is 0 -> 0, 1-5 -> n*40+55
			// https://github.com/gnachman/iTerm2/blob/5fc45c349417b8483dfe8426432fcbadc32cb6d9/sources/NSColor%2BiTerm.m#L335
			// Is this documented somewhere?
			f := func(c int) int {
				if c == 0 {
					return 0
				}
				return c*40 + 55
			}
			return Color{RGB: []int{f(r), f(g), f(b)}}, 2
		case n >= 232 && n <= 255:
			// 232-255:  grayscale from black to white in 24 steps
			g := int(255 * ((float32(n) - 232.0) / 23))
			return Color{RGB: []int{g, g, g}}, 2
		}
	}
	return Color{N: -1}, 0
}

var paramSplitRE = regexp.MustCompile(`[:;]`)

// ReadRune returns next rune. The decoder struct has state for last returned rune, .X, .Y, .Foreground etc.
func (d *Decoder) ReadRune() (r rune, size int, err error) {
	for {
		r, n, err := d.readBuf.ReadRune()
		if err != nil {
			return r, n, err
		}
		switch d.State {
		case StateCopy:
			switch r {
			case ESCRune:
				d.State = StateSeenESC
			default:
				d.X = d.nx
				d.Y = d.ny
				if d.Y > d.MaxY {
					d.MaxY = d.Y
				}

				switch r {
				case '\r':
					d.nx = 0
				case '\n':
					d.nx = 0
					d.ny++
				case '\t':
					d.nx += 8 - (d.nx % 8)
				default:
					if d.X > d.MaxX {
						d.MaxX = d.X
					}
					d.nx++
				}

				return r, n, err
			}
		case StateSeenESC:
			switch r {
			case '[':
				d.State = StateCSI
			case ']':
				d.State = StateOSC
			default:
				d.State = StateCopy
				return r, n, err
			}
		case StateCSI:
			switch {
			case bytes.ContainsAny([]byte(string([]rune{r})), FinalBytes):
				s := d.paramsBuf.String()
				ss := paramSplitRE.Split(s, -1)
				var pn []int
				for _, p := range ss {
					// will treat empty as 0
					n, _ := strconv.Atoi(p)
					pn = append(pn, n)
				}
				d.paramsBuf.Reset()

				switch r {
				case SGRByte:
					for i := 0; i < len(pn); i++ {
						n := pn[i]
						var ns int
						switch {
						case sgrReset.Is(n):
							d.Foreground = Color{N: -1}
							d.Background = Color{N: -1}
							d.Underline = false
							d.Intensity = false
							d.Invert = false
						case sgrIncreaseIntensity.Is(n):
							d.Intensity = true
						case sgrNormal.Is(n):
							d.Intensity = false
						case sgrForeground.Is(n):
							d.Foreground = Color{N: n - 30}
						case sgrForegroundBright.Is(n):
							d.Foreground = Color{N: n - 90 + 8}
						case sgrForegroundRGB.Is(n):
							d.Foreground, ns = intsToColor(30, 90, pn[i+1:])
							i += ns
						case sgrForegroundDefault.Is(n):
							d.Foreground = Color{N: -1}
						case sgrBackground.Is(n):
							d.Background = Color{N: n - 40}
						case sgrBackgroundBright.Is(n):
							d.Background = Color{N: n - 100 + 8}
						case sgrBackgroundRGB.Is(n):
							d.Background, ns = intsToColor(40, 100, pn[i+1:])
							i += ns
						case sgrBackgroundDefault.Is(n):
							d.Background = Color{N: -1}
						case sgrUnderlineOn.Is(n):
							d.Underline = true
						case sgrUnderlineOff.Is(n):
							d.Underline = false
						case sgrInvertOn.Is(n):
							d.Invert = true
						case sgrInvertOff.Is(n):
							d.Invert = false
						}
					}
				}
				d.State = StateCopy
			default:
				if _, err := d.paramsBuf.WriteRune(r); err != nil {
					return 0, 0, err
				}
			}
		case StateOSC:
			switch r {
			case BELRune:
				d.State = StateCopy
			case ESCRune:
				d.State = StateOSCSeenESC
			default:
				// nop, skip
			}
		case StateOSCSeenESC:
			switch r {
			case '\\':
				d.State = StateCopy
			default:
				// nop, skip
			}
		default:
			panic("unreachable")
		}
	}
}
