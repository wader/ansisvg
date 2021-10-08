package schemes

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/wader/ansisvg/colorscheme"
)

//go:embed *.json
var FS embed.FS

func Load(name string) (colorscheme.WorkbenchColorCustomizations, error) {
	var vsCS colorscheme.VSCodeColorScheme
	f, err := FS.Open(name + ".json")
	if err != nil {
		return vsCS.WorkbenchColorCustomizations, fmt.Errorf("scheme not found")
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&vsCS); err != nil {
		return vsCS.WorkbenchColorCustomizations, err
	}

	return vsCS.WorkbenchColorCustomizations, nil
}
