package esms

import (
	"fmt"
	"strings"
)

type PlayerName string

func (p PlayerName) String() string {
	return string(p)
}

func (p PlayerName) Short(max int) string {
	parts := strings.Split(string(p), " ")
	if len(parts) == 1 {
		return parts[0]
	}
	name := fmt.Sprintf("%s. %s", parts[0][0:1], strings.Join(parts[1:], " "))
	if len(name) > max {
		return name[0:max]
	}
	return name
}

type TeamPlayer struct {
	Player *Player
	Name   PlayerName

	Active  int // 1 if the player is currently active on the field, 0 - TODO: should be a bool
	Minutes int // number of minutes played

	Pos  string // position: "GK", "DF", "MF", "FW"
	Side string // "L", "C", "R" - left, center, right

	pref_side    string
	likes_right  bool
	likes_left   bool
	likes_center bool

	ag float64

	// Skills
	st int
	tk int
	ps int
	sh int

	// Abilities
	sh_ab int
	ps_ab int
	tk_ab int
	st_ab int

	fatigue                    float64
	nominal_fatigue_per_minute float64
	injured                    int
	fitness                    int
	stamina                    int

	injury     int
	suspension int

	// Contributions
	tk_contrib float64
	ps_contrib float64
	sh_contrib float64

	// Stats
	yellowcards int
	redcards    int
	goals       int
	assists     int
	tackles     int
	keypasses   int
	shots       int
	shots_on    int
	shots_off   int
	saves       int
	conceded    int
	fouls       int
}

type Team struct {
	Name          string
	Players       [20]TeamPlayer
	Score         int
	CurrentGK     int
	Substitutions int
	Aggression    float64
	FinalFouls    int
	PenaltyTaker  int // index of the player who takes penalties
	TeamTackling  float64
	TeamPassing   float64
	TeamShooting  float64
	ShotProb      float64
	FinalShotsOn  int
	FinalShotsOff int
	Tactic        string
	Injuries      int
	Roster        Roster
	Colors        []string
}

type Teamsheet struct {
	Tactic string            `json:"tactic"`
	Field  []TeamsheetPlayer `json:"field"`
	Bench  []TeamsheetPlayer `json:"bench"`
	PK     string            `json:"pk"`
	Name   string            `json:"team_name"`
}

type TeamsheetPlayer struct {
	Name PlayerName `json:"name"`
	Pos  string     `json:"pos"`
}

type Roster struct {
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
