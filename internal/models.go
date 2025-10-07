package internal

import (
	"strings"
)

type Person struct {
	FirstName string
	LastName  string
	Country   string
}

func (p Person) String() string {
	return strings.Join([]string{p.FirstName, p.LastName}, " ")
}

func (p Person) Short() string {
	fn := ""
	parts := strings.Split(p.FirstName, " ")
	for _, p := range parts {
		fn += strings.ToUpper(p[0:1])
	}
	return strings.Join([]string{fn, p.LastName}, " ")
}

type Roster struct {
	TeamCode string  `json:"team_code"`
	TeamName string  `json:"team_name"`
	Players  Players `json:"players"`
}

type Players []Player

func (p Players) FindByName(name string) *Player {
	for _, pl := range p {
		if pl.Name == name {
			return &pl
		}
	}
	return nil
}

type Player struct {
	ID          string
	Ag          float64 `json:"ag"`
	Age         int     `json:"age"`
	Assists     int     `json:"assists"`
	Dp          int     `json:"dp"`
	Fitness     int     `json:"fitness"`
	Games       int     `json:"games"`
	Goals       int     `json:"goals"`
	Injury      int     `json:"injury"`
	KeyPasses   int     `json:"keypasses"`
	Name        string  `json:"name"`
	Nationality string  `json:"nationality"`
	PrefSide    string  `json:"pref_side"`
	Ps          int     `json:"ps"`
	PsAb        int     `json:"ps_ab"`
	Saves       int     `json:"saves"`
	Sh          int     `json:"sh"`
	ShAb        int     `json:"sh_ab"`
	Shots       int     `json:"shots"`
	St          int     `json:"st"`
	StAb        int     `json:"st_ab"`
	Stamina     int     `json:"stamina"`
	Suspension  int     `json:"suspension"`
	Tackles     int     `json:"tackles"`
	Team        string  `json:"team"`
	Tk          int     `json:"tk"`
	TkAb        int     `json:"tk_ab"`
}
