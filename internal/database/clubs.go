package database

import (
	"path/filepath"

	"github.com/ilmaruk/go-esms/internal"
)

func LoadAllClubs(rootDir string) ([]internal.Club, error) {
	var clubs []internal.Club
	err := loadData(filepath.Join(dataDir(rootDir), "clubs.json"), &clubs)
	return clubs, err
}

func LoadClubsByLeague(rootDir, league string) ([]internal.Club, error) {
	allClubs, err := LoadAllClubs(rootDir)
	if err != nil {
		return nil, err
	}

	// Filter by league
	var clubs []internal.Club
	for _, c := range allClubs {
		if c.League == league {
			clubs = append(clubs, c)
		}
	}
	return clubs, err
}
