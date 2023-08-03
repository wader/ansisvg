// Package schemes embeds a database of color schemes from https://github.com/mbadolato/iTerm2-Color-Schemes
package schemes

import (
	"embed"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/wader/ansisvg/colorscheme"
)

//go:embed *.json
var fs embed.FS

const jsonExt = ".json"

func Load(name string) (colorscheme.WorkbenchColorCustomizations, error) {
	var vsCS colorscheme.VSCodeColorScheme
	f, err := fs.Open(name + jsonExt)
	if err != nil {
		return vsCS.WorkbenchColorCustomizations, fmt.Errorf("scheme not found")
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&vsCS); err != nil {
		return vsCS.WorkbenchColorCustomizations, err
	}

	return vsCS.WorkbenchColorCustomizations, nil
}

func Names() []string {
	es, err := fs.ReadDir(".")
	if err != nil {
		// should not happen
		panic(err)
	}

	var ns []string
	for _, e := range es {
		n := e.Name()
		ns = append(ns, e.Name()[0:len(n)-len(jsonExt)])
	}
	sort.Strings(ns)

	return ns
}
