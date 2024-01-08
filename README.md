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
ansisvg - Convert ANSI to SVG
Usage: ansisvg [FLAGS]

Example usage:
  program | ansisvg > file.svg

--charboxsize       Character box size (use pixel units instead of font units)
--colorscheme       Color scheme
--fontfile          Font file to use and embed
--fontname          Font name
--fontref           External font URL to use
--fontsize          Font size
--grid              Grid mode (sets position for each character)
--help, -h          Show help
--listcolorschemes  List color schemes
--transparent       Transparent background
--version, -v       Show version
--width, -w         Terminal width (auto if not set)
```

Color themes are the ones from https://github.com/mbadolato/iTerm2-Color-Schemes

## Install and build

To build you will need at least go 1.16 or later.

Install latest master and copy it to `/usr/local/bin`:
```sh
go install github.com/wader/ansisvg@master
cp $(go env GOPATH)/bin/ansisvg /usr/local/bin
```

Build from cloned repo:
```
go build -o ansisvg .
```

## Fonts

Note that embedded fonts might not be supported by some SVG viewers. At time of writing this is not supported by Inkscape (see https://gitlab.com/inkscape/inbox/-/issues/301).

## Tricks

### ANSI to PDF or PNG

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

## Development and release build

Run all tests and write new test output:
```
go test ./... -update
```

Manual release build with version can be done with:
```
go build -ldflags "-X main.version=1.2.3" -o ansisvg .
```

Visual inspect test output in browser:
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

 UbuntuMonoNerdFontMono-Regular.woff2 from https://github.com/ryanoasis/nerd-fonts license https://github.com/ryanoasis/nerd-fonts/blob/master/LICENSE

## TODO and ideas
- Underline overlaps a bit, sometimes causing weird blending
- Handle vertical tab and form feed (normalize into spaces?)
- Handle overdrawing
- More CSI, keep track of cursor?
- PNG output (embed nice fonts?)
