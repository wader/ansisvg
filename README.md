# ansisvg

Convert ANSI output to SVG.

## Usage

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
    	Character box size (default 7x13)
  -colorscheme string
    	Color scheme (default "Builtin Dark")
  -font string
    	Font (default "Monaco, Lucida Console, Courier")
  -fontsize int
    	Font size (default 12)
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

## Licenses and thanks

Color themes from
https://github.com/mbadolato/iTerm2-Color-Schemes,
license https://github.com/mbadolato/iTerm2-Color-Schemes/blob/master/LICENSE

colortest from https://github.com/pablopunk/colortest

## TODO

- Bold
- Foreground black with intensity
- More CSI, keep track of cursor?
