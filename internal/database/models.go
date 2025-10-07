package database

import "github.com/google/uuid"

type Player struct{}

type Roster struct {
	Name    string      `json:"team_name"`
	Players []uuid.UUID `json:"players"`
}
