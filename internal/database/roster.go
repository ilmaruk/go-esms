package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilmaruk/go-esms/internal"
)

func LoadRoster(rootDir, code string) (internal.Roster, error) {
	path := filepath.Join(dataDir(rootDir), "teams", fmt.Sprintf("%s.json", code))
	var r internal.Roster
	err := loadData(path, &r)
	return r, err
}

func SaveRoster(rootDir string, roster internal.Roster) error {
	b, err := json.MarshalIndent(roster, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(dataDir(rootDir), "teams", strings.ToLower(fmt.Sprintf("%s.json", roster.TeamCode)))
	return os.WriteFile(path, b, 0644)
}
