package ansitosvg_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/wader/ansisvg/ansitosvg"
	"github.com/wader/ansisvg/internal/difftest"
)

func TestCovert(t *testing.T) {
	difftest.TestWithOptions(t, difftest.Options{
		Path:        "testdata",
		Pattern:     "*.ansi",
		ColorDiff:   os.Getenv("TEST_COLOR") != "",
		WriteOutput: os.Getenv("WRITE_ACTUAL") != "",
		Fn: func(t *testing.T, path, input string) (string, string, error) {
			opts := ansitosvg.DefaultOptions
			optsPath := path + ".json"
			if f, err := os.Open(optsPath); err == nil {
				defer f.Close()
				if err := json.NewDecoder(f).Decode(&opts); err != nil {
					t.Fatal(err)
				}
			}

			actual := &bytes.Buffer{}
			err := ansitosvg.Convert(
				bytes.NewBufferString(input),
				actual,
				opts,
			)
			if err != nil {
				return "", "", err
			}

			return path + ".svg", actual.String(), nil
		},
	})
}
