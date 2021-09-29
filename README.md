# ansisvg

Convert ANSI terminal codes to SVG.

```sh
./colortest | ansivg > doc/example.svg
 ```
![doc/colortest.svg asdad](doc/colortest.svg)

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
