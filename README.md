# ansisvg

Convert ANSI output to SVG.

Pipe output from program thru ansisvg and it will output a SVG file on stdout.
```sh
./colortest | ansisvg > doc/example.svg
 ```
Produces [colortest.svg](doc/colortest.svg)

![doc/colortest.svg asdad](doc/colortest.svg)

```
$ ansisvg -h
Usage of ansisvg:
  -charboxsize value
    	Character box size (default 8x16)
  -colorscheme string
    	Color scheme (default "Builtin Dark")
  -fontname string
    	Font name (default "Monaco, Lucida Console, Courier")
  -fontsize int
    	Font size (default 14)
  -transparent
    	Transparent background
  -width int
    	Terminal width (auto)
```

Color themes are the ones from https://github.com/mbadolato/iTerm2-Color-Schemes

## Install

```sh
# build from cloned repo
go build -o ansisvg main.go

# install directly
go install github.com/wader/ansisvg@master
# copy binary to $PATH
cp $(go env GOPATH)/bin/ansisvg /usr/local/bin
```

## Tricks to force color output from programs

### script
`script -q /dev/null <command> | ansisvg`
### ffmpeg
`TERM=a AV_LOG_FORCE_COLOR=1 ffmpeg ... 2>&1 | ansisvg`
### jq
`jq -C | ansisvg`

## Licenses and thanks

Color themes from
https://github.com/mbadolato/iTerm2-Color-Schemes,
license https://github.com/mbadolato/iTerm2-Color-Schemes/blob/master/LICENSE

colortest from https://github.com/pablopunk/colortest

## TODO and ideas
- Bold
- More CSI, keep track of cursor?
- PNG output
