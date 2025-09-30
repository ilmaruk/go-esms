package esms

type Player struct {
	Active  int // 1 if the player is currently active on the field, 0 - TODO: should be a bool
	Minutes int // number of minutes played

	Pos  string // position: "GK", "DF", "MF", "FW"
	Side string // "L", "C", "R" - left, center, right

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
	Players       [numPlayers]Player
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
}
