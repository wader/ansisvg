// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package getopt parses command lines using getopt(3) syntax.
// It is a replacement for flag.Parse but still expects flags themselves
// to be defined in package flag.
//
// Flags defined with one-letter names are available as short flags
// (invoked using one dash, as in -x) and all flags are available as
// long flags (invoked using two dashes, as in --x or --xylophone).
//
// To use, define flags as usual with package flag. Then introduce
// any aliases by calling getopt.Alias:
//
//	getopt.Alias("n", "dry-run")
//	getopt.Alias("v", "verbose")
//
// Or call getopt.Aliases to define a list of aliases:
//
//	getopt.Aliases(
//		"n", "dry-run",
//		"v", "verbose",
//	)
//
// One name in each pair must already be defined in package flag
// (so either "n" or "dry-run", and also either "v" or "verbose").
//
// Then parse the command-line:
//
//	getopt.Parse()
//
// If it encounters an error, Parse calls flag.Usage and then exits the program.
//
// When writing a custom flag.Usage function, call getopt.PrintDefaults
// instead of flag.PrintDefaults to get a usage message that includes the
// names of aliases in flag descriptions.
//
// At initialization time, this package installs a new flag.Usage that is the
// same as the default flag.Usage except that it calls getopt.PrintDefaults
// instead of flag.PrintDefaults.
//
// This package also defines a FlagSet wrapping the standard flag.FlagSet.
//
// # Caveat
//
// In general Go flag parsing is preferred for new programs, because
// it is not as pedantic about the number of dashes used to invoke
// a flag (you can write -verbose or --verbose and the program
// does not care). This package is meant to be used in situations
// where, for legacy reasons, it is important to use exactly getopt(3)
// syntax, such as when rewriting in Go an existing tool that already
// uses getopt(3).
package getopt // import "rsc.io/getopt"

import (
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"unicode/utf8"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		PrintDefaults() // ours not package flag's
	}

	CommandLine.FlagSet = flag.CommandLine
	CommandLine.name = os.Args[0]
	CommandLine.errorHandling = flag.ExitOnError
	CommandLine.outw = os.Stderr
	CommandLine.Usage = func() { flag.Usage() }
}

var CommandLine FlagSet

// A FlagSet is a set of defined flags.
// It wraps and provides the same interface as flag.FlagSet
// but parses command line arguments using getopt syntax.
//
// Note that "go doc" shows only the methods customized
// by package getopt; FlagSet also provides all the methods
// of the embedded flag.FlagSet, like Bool, Int, NArg, and so on.
type FlagSet struct {
	*flag.FlagSet

	alias         map[string]string
	unalias       map[string]string
	name          string
	errorHandling flag.ErrorHandling
	outw          io.Writer
}

func (f *FlagSet) out() io.Writer {
	if f.outw == nil {
		return os.Stderr
	}
	return f.outw
}

// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
func (f *FlagSet) SetOutput(output io.Writer) {
	f.FlagSet.SetOutput(output)
	f.outw = output
}

// NewFlagSet returns a new, empty flag set with the specified name and error
// handling property.
func NewFlagSet(name string, errorHandling flag.ErrorHandling) *FlagSet {
	f := new(FlagSet)
	f.Init(name, errorHandling)
	return f
}

// Init sets the name and error handling proprety for a flag set.
func (f *FlagSet) Init(name string, errorHandling flag.ErrorHandling) {
	if f.FlagSet == nil {
		f.FlagSet = new(flag.FlagSet)
	}
	f.FlagSet.Init(name, errorHandling)
	f.name = name
	f.errorHandling = errorHandling
	f.FlagSet.Usage = f.defaultUsage
}

func (f *FlagSet) init() {
	if f.alias == nil {
		f.alias = make(map[string]string)
		f.unalias = make(map[string]string)
	}
}

func (f *FlagSet) LookupShort(name string) string {
	return f.unalias[name]
}

// Lookup returns the Flag structure of the named flag,
// returning nil if none exists.
// If name is a defined alias for a defined flag,
// Lookup returns the original flag; in this case
// the Name field in the result will differ from the
// name passed to Lookup.
func (f *FlagSet) Lookup(name string) *flag.Flag {
	if x, ok := f.alias[name]; ok {
		name = x
	}
	return f.FlagSet.Lookup(name)
}

// Alias introduces an alias for an existing flag name.
// The short name must be a single letter, and the long name must be multiple letters.
// Exactly one name must be defined as a flag already: the undefined name is introduced
// as an alias for the defined name.
// Alias panics if both names are already defined or if both are undefined.
//
// For example, if a flag named "v" is already defined using package flag,
// then it is available as -v (or --v). Calling Alias("v", "verbose") makes the same
// flag also available as --verbose.
func Alias(short, long string) {
	CommandLine.Alias(short, long)
}

// Alias introduces an alias for an existing flag name.
// The short name must be a single letter, and the long name must be multiple letters.
// Exactly one name must be defined as a flag already: the undefined name is introduced
// as an alias for the defined name.
// Alias panics if both names are already defined or if both are undefined.
//
// For example, if a flag named "v" is already defined using package flag,
// then it is available as -v (or --v). Calling Alias("v", "verbose") makes the same
// flag also available as --verbose.
func (f *FlagSet) Alias(short, long string) {
	f.init()
	if short == "" || long == "" {
		panic("Alias: invalid empty flag name")
	}
	if utf8.RuneCountInString(short) != 1 {
		panic("Alias: invalid short flag name -" + short)
	}
	if utf8.RuneCountInString(long) == 1 {
		panic("Alias: invalid long flag name --" + long)
	}

	f1 := f.Lookup(short)
	f2 := f.Lookup(long)
	if f1 == nil && f2 == nil {
		panic("Alias: neither -" + short + " nor -" + long + " is a defined flag")
	}
	if f1 != nil && f2 != nil {
		panic("Alias: both -" + short + " and -" + long + " are defined flags")
	}

	if f1 != nil {
		f.alias[long] = short
		f.unalias[short] = long
	} else {
		f.alias[short] = long
		f.unalias[long] = short
	}
}

// Aliases introduces zero or more aliases. The argument list must consist of an
// even number of strings making up a sequence of short, long pairs to be passed
// to Alias.
func Aliases(list ...string) {
	CommandLine.Aliases(list...)
}

// Aliases introduces zero or more aliases. The argument list must consist of an
// even number of strings making up a sequence of short, long pairs to be passed
// to Alias.
func (f *FlagSet) Aliases(list ...string) {
	if len(list)%2 != 0 {
		panic("getopt: Aliases not invoked with pairs")
	}
	for i := 0; i < len(list); i += 2 {
		f.Alias(list[i], list[i+1])
	}
}

type boolFlag interface {
	IsBoolFlag() bool
}

func (f *FlagSet) failf(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	fmt.Fprintln(f.out(), err)
	f.Usage()
	return err
}

// defaultUsage is the default function to print a usage message.
func (f *FlagSet) defaultUsage() {
	if f.name == "" {
		fmt.Fprintf(f.out(), "Usage:\n")
	} else {
		fmt.Fprintf(f.out(), "Usage of %s:\n", f.name)
	}
	f.PrintDefaults()
}

// Parse parses the command-line flags from os.Args[1:].
func Parse() {
	CommandLine.Parse(os.Args[1:])
}

// Parse parses flag definitions from the argument list,
// which should not include the command name.
// Parse must be called after all flags and aliases in the FlagSet are defined
// and before flags are accessed by the program.
// The return value will be flag.ErrHelp if -h or --help were used but not defined.
func (f *FlagSet) Parse(args []string) error {
	for len(args) > 0 {
		arg := args[0]
		if len(arg) < 2 || arg[0] != '-' {
			break
		}
		args = args[1:]
		if arg[:2] == "--" {
			// Process single long option.
			if arg == "--" {
				break
			}
			name := arg[2:]
			value := ""
			haveValue := false
			if i := strings.Index(name, "="); i >= 0 {
				name, value = name[:i], name[i+1:]
				haveValue = true
			}
			fg := f.Lookup(name)
			if fg == nil {
				if name == "h" || name == "help" {
					// TODO ErrHelp
				}
				return f.failf("flag provided but not defined: --%s", name)
			}
			if b, ok := fg.Value.(boolFlag); ok && b.IsBoolFlag() {
				if haveValue {
					if err := fg.Value.Set(value); err != nil {
						return f.failf("invalid boolean value %q for --%s: %v", value, name, err)
					}
				} else {
					if err := fg.Value.Set("true"); err != nil {
						return f.failf("invalid boolean flag %s: %v", name, err)
					}
				}
				continue
			}
			if !haveValue {
				if len(args) == 0 {
					return f.failf("missing argument for --%s", name)
				}
				value, args = args[0], args[1:]
			}
			if err := fg.Value.Set(value); err != nil {
				return f.failf("invalid value %q for flag --%s: %v", value, name, err)
			}
			continue
		}

		// Process one or more short options.
		for arg = arg[1:]; arg != ""; {
			r, size := utf8.DecodeRuneInString(arg)
			if r == utf8.RuneError && size == 1 {
				return f.failf("invalid UTF8 in command-line flags")
			}
			name := arg[:size]
			arg = arg[size:]
			fg := f.Lookup(name)
			if fg == nil {
				if name == "h" {
					// TODO ErrHelp
				}
				return f.failf("flag provided but not defined: -%s", name)
			}
			if b, ok := fg.Value.(boolFlag); ok && b.IsBoolFlag() {
				if err := fg.Value.Set("true"); err != nil {
					return f.failf("invalid boolean flag %s: %v", name, err)
				}
				continue
			}
			if arg == "" {
				if len(args) == 0 {
					return f.failf("missing argument for -%s", name)
				}
				arg, args = args[0], args[1:]
			}
			if err := fg.Value.Set(arg); err != nil {
				return f.failf("invalid value %q for flag -%s: %v", arg, name, err)
			}
			break // consumed arg
		}
	}

	// Arrange for flag.NArg, flag.Args, etc to work properly.
	f.FlagSet.Parse(append([]string{"--"}, args...))
	return nil
}

// PrintDefaults is like flag.PrintDefaults but includes information
// about short/long alias pairs and prints the correct syntax for
// long flags.
func PrintDefaults() {
	CommandLine.PrintDefaults()
}

// PrintDefaults is like flag.PrintDefaults but includes information
// about short/long alias pairs and prints the correct syntax for
// long flags.
func (f *FlagSet) PrintDefaults() {
	f.FlagSet.VisitAll(func(fg *flag.Flag) {
		name := fg.Name
		short, long := "", ""
		other := f.unalias[name]
		if utf8.RuneCountInString(name) > 1 {
			long, short = name, other
		} else {
			short, long = name, other
		}
		var s string
		if short != "" {
			s = fmt.Sprintf("  -%s", short) // Two spaces before -; see next two comments.
			if long != "" {
				s += ", --" + long
			}
		} else {
			s = fmt.Sprintf("  --%s", long) // Two spaces before -; see next two comments.
		}
		name, usage := flag.UnquoteUsage(fg)
		if len(name) > 0 {
			s += " " + name
		}

		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if len(s) <= 4 { // space, space, '-', 'x'.
			s += "\t"
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			s += "\n    \t"
		}
		s += usage
		if !isZeroValue(fg, fg.DefValue) {
			if strings.HasSuffix(reflect.TypeOf(fg.Value).String(), "stringValue") {
				// put quotes on the value
				s += fmt.Sprintf(" (default %q)", fg.DefValue)
			} else {
				s += fmt.Sprintf(" (default %v)", fg.DefValue)
			}
		}
		fmt.Fprint(f.out(), s, "\n")
	})
}

// isZeroValue guesses whether the string represents the zero
// value for a flag. It is not accurate but in practice works OK.
func isZeroValue(f *flag.Flag, value string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(f.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	if value == z.Interface().(flag.Value).String() {
		return true
	}

	switch value {
	case "false":
		return true
	case "":
		return true
	case "0":
		return true
	}
	return false
}
