package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilmaruk/go-esms/internal"
)

func SaveTeamsheet(workDir string, teamsheet internal.Teamsheet) error {
	b, err := json.MarshalIndent(teamsheet, "", " ")
	if err != nil {
		return err
	}
	path := filepath.Join(workDir, "data", "teamsheets", strings.ToLower(fmt.Sprintf("%s.json", teamsheet.Code)))
	return os.WriteFile(path, b, 0644)
}
