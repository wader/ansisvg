package ansisvg_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/wader/ansisvg/ansisvg"
	"github.com/wader/ansisvg/internal/deepequal"
)

func TestCovert(t *testing.T) {
	err := filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".ansi" {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			ansiInput, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			expectedSVG, err := ioutil.ReadFile(path + ".svg")
			if err != nil {
				t.Fatal(err)
			}

			actualSVG := &bytes.Buffer{}
			err = ansisvg.Convert(
				bytes.NewReader(ansiInput),
				actualSVG,
				ansisvg.DefaultOptions,
			)
			if err != nil {
				t.Fatal(err)
			}

			deepequal.Error(t, "svg", string(expectedSVG), actualSVG.String())
		})

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
