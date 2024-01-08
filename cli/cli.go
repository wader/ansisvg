package cli

import (
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/wader/ansisvg/ansitosvg"
	"github.com/wader/ansisvg/colorscheme/schemes"
	"github.com/wader/ansisvg/getopt"
	"github.com/wader/ansisvg/svgscreen"
)

type boxSize struct {
	Width  int
	Height int
}

func (d *boxSize) String() string {
	return fmt.Sprintf("%dx%d", d.Width, d.Height)
}

func (d *boxSize) Set(s string) error {
	parts := strings.Split(s, "x")
	if len(parts) != 2 {
		return fmt.Errorf("must be WxH")
	}
	d.Width, _ = strconv.Atoi(parts[0])
	d.Height, _ = strconv.Atoi(parts[1])
	return nil
}

type Env struct {
	ReadFile func(string) ([]byte, error)
	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
	Args     []string
}

func Main(env Env) error {
	fs := getopt.NewFlagSet("ansisvg", flag.ContinueOnError)
	var fontNameFlag = fs.String("fontname", ansitosvg.DefaultOptions.FontName, "Font name")
	var fontFileFlag = fs.String("fontfile", "", "Font file to use and embed")
	var fontRefFlag = fs.String("fontref", "", "External font file to reference")
	var fontSizeFlag = fs.Int("fontsize", ansitosvg.DefaultOptions.FontSize, "Font size")
	var terminalWidthFlag = fs.Int("width", 0, "Terminal width (auto)")
	var colorSchemeFlag = fs.String("colorscheme", ansitosvg.DefaultOptions.ColorScheme, "Color scheme")
	var listColorSchemesFlag = fs.Bool("listcolorschemes", false, "List color schemes")
	var transparentFlag = fs.Bool("transparent", ansitosvg.DefaultOptions.Transparent, "Transparent background")
	var gridModeFlag = fs.Bool("grid", false, "Enable grid mode (sets position for each character)")
	var helpFlag = fs.Bool("help", false, "Show help")
	var charBoxSize = boxSize{
		Width:  ansitosvg.DefaultOptions.CharBoxSize.Width,
		Height: ansitosvg.DefaultOptions.CharBoxSize.Height,
	}
	fs.Var(&charBoxSize, "charboxsize", "Character box size (forces pixel units instead of font-relative units)")
	// handle error and usage output ourself
	fs.Usage = func() {}
	fs.SetOutput(io.Discard)
	usage := func() {
		maxNameLen := 0
		fs.VisitAll(func(f *flag.Flag) {
			if len(f.Name) > maxNameLen {
				maxNameLen = len(f.Name)
			}
		})

		fmt.Fprintf(env.Stdout, `
%[1]s - Convert ANSI to SVG
Usage: %[1]s [FLAGS]
`[1:], fs.Name())
		fs.VisitAll(func(f *flag.Flag) {
			def := ""
			valueStr := f.Value.String()
			if valueStr != "" && valueStr != "true" && valueStr != "false" && valueStr != "0" {
				def = " (" + valueStr + ")"
			}
			short := ""
			if s := fs.LookupShort(f.Name); s != "" {
				short = ", -" + s
			}

			flagNames := f.Name + short
			pad := strings.Repeat(" ", maxNameLen-len(flagNames))
			fmt.Fprintf(env.Stdout, "--%s%s%s  %s%s\n", f.Name, short, pad, f.Usage, def)
		})
	}
	fs.Aliases(
		"h", "help",
	)
	if err := fs.Parse(env.Args[1:]); err != nil {
		return err
	}
	if *helpFlag {
		usage()
		return nil
	}

	if *listColorSchemesFlag {
		maxNameLen := 0
		for _, n := range schemes.Names() {
			if len(n) > maxNameLen {
				maxNameLen = len(n)
			}
		}
		for _, n := range schemes.Names() {
			s, _ := schemes.Load(n)
			pad := strings.Repeat(" ", maxNameLen+1-len(n))
			fmt.Fprintf(env.Stdout, "%s\n", s.ANSIDemo(n+pad))
		}
		return nil
	}

	var fontEmbedded []byte
	if *fontFileFlag != "" {
		var err error
		fontEmbedded, err = env.ReadFile(*fontFileFlag)
		if err != nil {
			return err
		}
	}

	return ansitosvg.Convert(
		env.Stdin,
		env.Stdout,
		ansitosvg.Options{
			FontName:      *fontNameFlag,
			FontEmbedded:  fontEmbedded,
			FontRef:       *fontRefFlag,
			FontSize:      *fontSizeFlag,
			TerminalWidth: *terminalWidthFlag,
			CharBoxSize: svgscreen.BoxSize{
				Width:  charBoxSize.Width,
				Height: charBoxSize.Height,
			},
			ColorScheme: *colorSchemeFlag,
			Transparent: *transparentFlag,
			GridMode:    *gridModeFlag,
		},
	)
}
