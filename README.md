# ansisvg

Convert [ANSI](https://en.wikipedia.org/wiki/ANSI_escape_code) to [SVG](https://en.wikipedia.org/wiki/Scalable_Vector_Graphics).

Pipe output from a program thru ansisvg and it will output a SVG file on stdout. Can be used to
produce nice looking example output for presentations, markdown files etc.

```sh
./colortest | ansisvg > colortest.ansi.svg
 ```
Produces [colortest.ansi.svg](ansitosvg/testdata/colortest.ansi.svg)

![ansisvg output for colortest](ansitosvg/testdata/colortest.ansi.svg)

```
$ ansisvg -h
Usage of ansisvg:
  -charboxsize value
    	Character box size (default 8x16)
  -colorscheme string
    	Color scheme (default "Builtin Dark")
  -fontname string
    	Font name
  -fontsize int
    	Font size (default 14)
  -transparent
    	Transparent background
  -width int
    	Terminal width (auto)
```

Color themes are the ones from https://github.com/mbadolato/iTerm2-Color-Schemes

## Install

Install latest master and copy it to `/usr/local/bin`:
```sh
go install github.com/wader/ansisvg@master
cp $(go env GOPATH)/bin/ansisvg /usr/local/bin
```

## Development

Build from cloned repo:
```
go build -o ansisvg main.go
```

## Tricks

### Use `bat` to produce source code highlighting
`bat --color=always -p main.go | ansisvg`

### Use `script` to run with a pty
`script -q /dev/null <command> | ansisvg`

### ffmpeg
`TERM=a AV_LOG_FORCE_COLOR=1 ffmpeg ... 2>&1 | ansisvg`

### jq
`jq -C | ansisvg`

## Development

Run all tests and write new difftest outputs
```
WRITE_ACTUAL=1 go test ./...
```

Visual inspect outputs in browser:
```
for i in ansisvg/testdata/*.ansi.svg; do echo "$i<br><img src=\"$i\"/><br>" ; done  > all.html
open all.html
```

Using [ffcat](https://github.com/wader/ffcat):
```
for i in ansisvg/testdata/*.ansi; do echo $i ; cat $i | go run main.go | inkscape --pipe --export-type=png -o - 2>/dev/null | ffcat ; done
```

## Licenses and thanks

Color themes from
https://github.com/mbadolato/iTerm2-Color-Schemes,
license https://github.com/mbadolato/iTerm2-Color-Schemes/blob/master/LICENSE

Uses colortest from https://github.com/pablopunk/colortest and terminal-colors from https://github.com/eikenb/terminal-colors.

## TODO and ideas
- Bold
- Underline overlaps a bit, sometimes causing weird blending
- Somehow use `<tspan>`/`textLength` to produce smaller output. Maybe `em/ch` CSS units for background rects,
but seems inkscape do not like `ch`. Would also make it nicer to copy text from SVG.
- Handle vertical tab and form feed (normalize into spaces?)
- Handle overdrawing
- More CSI, keep track of cursor?
- PNG output (embed nice fonts?)
