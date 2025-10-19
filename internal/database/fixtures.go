package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilmaruk/go-esms/internal"
)

func SaveFixtrues(rootDir string, fixtures internal.Fixtures) error {
	b, err := json.MarshalIndent(fixtures, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(dataDir(rootDir), "seasons", fmt.Sprintf("%03d", fixtures.Season), fmt.Sprintf("fixtures_%s.json", fixtures.League))
	return os.WriteFile(path, b, 0644)
}
