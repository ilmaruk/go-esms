package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilmaruk/go-esms/internal"
)

type DatabaseRepo struct {
	rootDir string
}

func NewDatabaseRepo(rootDir string) *DatabaseRepo {
	return &DatabaseRepo{rootDir: rootDir}
}

func (d *DatabaseRepo) Load(season uint, league string) (internal.Table, error) {
	path := filepath.Join(d.rootDir, "data", "seasons", fmt.Sprintf("%d", season), fmt.Sprintf("table_%s.json", league))
	var table internal.Table
	err := loadData(path, &table)
	return table, err
}

func SaveTable(rootDir string, table internal.Table, season int, league string) error {
	b, err := json.MarshalIndent(table, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(rootDir, "data", "seasons", fmt.Sprintf("%d", season), fmt.Sprintf("table_%s.json", league))
	return os.WriteFile(path, b, 0644)
}
