package database

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func loadData(path string, data any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, data)
}

func dataDir(rootDir string) string {
	return filepath.Join(rootDir, "data")
}
