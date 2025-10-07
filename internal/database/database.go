package database

import (
	"encoding/json"
	"os"
)

func loadData(path string, data any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, data)
}

// func loadTeamsheet(path string, ts *Teamsheet) error {
// 	return loadData(path, ts)
// }
