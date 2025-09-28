package esms

type Player struct {
	Active  int // 1 if the player is currently active on the field, 0 - TODO: should be a bool
	Minutes int // number of minutes played

	Pos string // position: "GK", "DF", "MF", "FW"

	// Abilities
	sh_ab int
	ps_ab int
	tk_ab int
	st_ab int

	// Stats
	yellowcards int
	redcards    int
	goals       int
	assists     int
	tackles     int
	keypasses   int
	shots_on    int
	shots_off   int
	saves       int
	conceded    int
}

type Team struct {
	Players   [numPlayers]Player
	Score     int
	CurrentGK int
}
