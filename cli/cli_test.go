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

func testHelper(t *testing.T, pattern string, ext string) {
	difftest.TestWithOptions(t, difftest.Options{
		Path:        "testdata",
		Pattern:     pattern,
		ColorDiff:   os.Getenv("TEST_COLOR") != "",
		WriteOutput: *update,
		Fn: func(t *testing.T, path, input string) (string, string, error) {
			testDir := filepath.Dir(path)
			testBase := filepath.Base(path)
			testName := testBase[0 : len(testBase)-len(filepath.Ext(testBase))]
			readFile := func(s string) ([]byte, error) {
				return os.ReadFile(filepath.Join(testDir, s))
			}
			readFileOrEmpty := func(s string) []byte {
				b, err := readFile(s)
				if err != nil {
					return nil
				}
				return b
			}

			args := string(readFileOrEmpty(testName + ".args"))
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

			return filepath.Join(testDir, testName) + ext, actualStdout.String(), nil
		},
	})
}

func TestMain(t *testing.T) {
	testHelper(t, "*.ansi", ".svg")
	testHelper(t, "*.stdin", ".stdout")
}
