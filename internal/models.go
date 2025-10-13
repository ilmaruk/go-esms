package internal

import (
	"strings"

	"github.com/google/uuid"
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

type Teamsheet struct {
	Tactic       string            `json:"tactic"`
	Field        []TeamsheetPlayer `json:"field"`
	Bench        []TeamsheetPlayer `json:"bench"`
	PenaltyTaker string            `json:"pk"`
	Name         string            `json:"team_name"`
	Code         string            `json:"team_code"`
}

type TeamsheetPlayer struct {
	Name string `json:"name"`
	Pos  string `json:"pos"`
}

type ClubColors struct {
	Primary   string `json:"primary"`
	Secondary string `jsson:"secondary"`
	Accent    string `json:"accent"`
}

type Club struct {
	Name   string     `json:"name"`
	Colors ClubColors `json:"colors"`
	Code   string     `json:"code"`
	City   string     `json:"city"`
	Elo    int        `json:"elo"`
	League string     `json:"league"`
}

type Game struct {
	ID   uuid.UUID `json:"ID"`
	Home string    `json:"home"`
	Away string    `json:"away"`
	Seed int64     `json:"seed"`
}

type GameWeek struct {
	ID    int    `json:"id"`
	Games []Game `json:"games"`
}

type Fixtures struct {
	League string          `json:"league"`
	Season int             `json:"season"`
	Clubs  map[string]Club `json:"clubs"`
	Weeks  []GameWeek      `json:"weeks"`
}
