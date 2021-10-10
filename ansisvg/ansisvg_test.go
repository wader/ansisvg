package ansisvg_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/wader/ansisvg/ansisvg"
	"github.com/wader/ansisvg/internal/difftest"
)

func TestCovert(t *testing.T) {
	difftest.TestWithOptions(t, difftest.Options{
		Path:        "testdata",
		Pattern:     "*.ansi",
		ColorDiff:   os.Getenv("TEST_COLOR") != "",
		WriteOutput: os.Getenv("WRITE_ACTUAL") != "",
		Fn: func(t *testing.T, path, input string) (string, string, error) {
			actual := &bytes.Buffer{}
			err := ansisvg.Convert(
				bytes.NewBufferString(input),
				actual,
				ansisvg.DefaultOptions,
			)
			if err != nil {
				return "", "", err
			}
			return path + ".svg", actual.String(), nil
		},
	})
}
