package schemes

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/wader/ansisvg/colorscheme"
)

//go:embed *.json
var FS embed.FS

func Load(name string) (colorscheme.VSCodeColorScheme, error) {
	var c colorscheme.VSCodeColorScheme
	f, err := FS.Open(name + ".json")
	if err != nil {
		return c, fmt.Errorf("scheme not found")
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return c, err
	}

	return c, nil
}
