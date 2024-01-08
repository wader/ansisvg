# ansisvg

Convert [ANSI](https://en.wikipedia.org/wiki/ANSI_escape_code) to [SVG](https://en.wikipedia.org/wiki/Scalable_Vector_Graphics).

Pipe output from a program thru `ansisvg` and it will output a SVG file on stdout.

Can be used to produce nice looking example output for presentations, markdown files etc. Note that it
does not support programs that do cursor movements like ncurses programs etc.

```sh
./colortest | ansisvg > colortest.ansi.svg
 ```
Produces [colortest.ansi.svg](cli/testdata/colortest.ansi.svg):

![ansisvg output for colortest](cli/testdata/colortest.ansi.svg)

```
$ ansisvg -h
Usage of ansisvg:
  -charboxsize value
    	Character box size (forces pixel units instead of font-relative units)
  -colorscheme string
    	Color scheme (default "Builtin Dark")
  -fontfile string
    	Font file to use and embed
  -fontname string
    	Font name (default "Courier")
  -fontref string
    	External font file to reference
  -fontsize int
    	Font size (default 14)
  -grid
    	Enable grid mode (sets position for each character)
  -listcolorschemes
    	List color schemes
  -transparent
    	Transparent background
  -width int
    	Terminal width (auto)
```

Color themes are the ones from https://github.com/mbadolato/iTerm2-Color-Schemes

## Install

To build you will need at least go 1.16 or later.

Install latest master and copy it to `/usr/local/bin`:
```sh
go install github.com/wader/ansisvg@master
cp $(go env GOPATH)/bin/ansisvg /usr/local/bin
```

## Fonts

Note that embedded fonts might not be supported by some SVG viewers. At time of writing this is not supported by Inkscape (see https://gitlab.com/inkscape/inbox/-/issues/301).

## Development

Build from cloned repo:
```
go build -o ansisvg main.go
```

## Tricks

#### ANSI to PDF or PNG

```
... | ansisvg | inkscape --pipe --export-type=pdf -o file.pdf
... | ansisvg | inkscape --pipe --export-type=png -o file.png
```


### Use `bat` to produce source code highlighting

```
bat --color=always -p main.go | ansisvg
```

### Use `script` to run with a pty

```
script -q /dev/null <command> | ansisvg
```

### ffmpeg

```
TERM=a AV_LOG_FORCE_COLOR=1 ffmpeg ... 2>&1 | ansisvg
```

### jq
```
jq -C | ansisvg
```

## Development

Run all tests and write new difftest outputs
```
WRITE_ACTUAL=1 go test ./...
```

Visual inspect outputs in browser:
```
for i in cli/testdata/*.ansi.svg; do echo "$i<br><img src=\"$i\"/><br>" ; done  > all.html
open all.html
```

Using [ffcat](https://github.com/wader/ffcat):
```
for i in cli/testdata/*.ansi; do echo $i ; cat $i | go run main.go | ffcat ; done
```

## Thanks

- Patrick Huesmann [@patrislav1](https://github.com/patrislav1) for better ANSI support and lots SVG output improvements.

## Licenses and thanks

Color themes from
https://github.com/mbadolato/iTerm2-Color-Schemes,
license https://github.com/mbadolato/iTerm2-Color-Schemes/blob/master/LICENSE

Uses colortest from https://github.com/pablopunk/colortest and terminal-colors from https://github.com/eikenb/terminal-colors.

## TODO and ideas
- Underline overlaps a bit, sometimes causing weird blending
- Handle vertical tab and form feed (normalize into spaces?)
- Handle overdrawing
- More CSI, keep track of cursor?
- PNG output (embed nice fonts?)
