package cli_test

import (
	"bytes"
	"encoding/csv"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wader/ansisvg/cli"
	"github.com/wader/ansisvg/internal/difftest"
)

var update = flag.Bool("update", false, "Update tests")

func argsSplit(s string) []string {
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = ' '
	rs, _ := r.Read()
	return rs
}

func TestMain(t *testing.T) {
	difftest.TestWithOptions(t, difftest.Options{
		Path:        "testdata",
		Pattern:     "*.ansi",
		ColorDiff:   os.Getenv("TEST_COLOR") != "",
		WriteOutput: *update || os.Getenv("WRITE_ACTUAL") != "",
		Fn: func(t *testing.T, path, input string) (string, string, error) {
			testBaseDir := filepath.Dir(path)
			readFile := func(s string) ([]byte, error) {
				return os.ReadFile(filepath.Join(testBaseDir, s))
			}
			readFileOrEmpty := func(s string) []byte {
				b, err := readFile(s)
				if err != nil {
					return nil
				}
				return b
			}

			args := string(readFileOrEmpty(filepath.Base(path) + ".args"))
			actualStdout := &bytes.Buffer{}
			actualStderr := &bytes.Buffer{}
			if err := cli.Main(cli.Env{
				ReadFile: readFile,
				Stdin:    strings.NewReader(input),
				Stdout:   actualStdout,
				Stderr:   actualStderr,
				Args:     append([]string{"ansisvg"}, argsSplit(string(args))...),
			}); err != nil {
				t.Error(err)
			}

			return path + ".svg", actualStdout.String(), nil
		},
	})
}
