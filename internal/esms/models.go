package esms

type Player struct {
	Name string

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
	Players       [20]Player
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
	RosterPlayer  []RosterPlayer
}

type Teamsheet struct {
	Tactic string   `json:"tactic"`
	Field  []Player `json:"field"`
	Bench  []Player `json:"bench"`
	PK     string   `json:"pk"`
	Name   string   `json:"team_name"`
}

type Roster struct {
	TeamName string         `json:"team_name"`
	Players  []RosterPlayer `json:"players"`
}

type RosterPlayer struct {
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
