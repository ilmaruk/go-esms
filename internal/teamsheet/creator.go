package teamsheet

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/ilmaruk/go-esms/internal"
)

var rnd *rand.Rand

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixMicro()))
}

func CreateTeamsheet(roster internal.Roster, tactic string) (internal.Teamsheet, error) {
	num_subs := 7

	t_player := make([]internal.TeamsheetPlayer, 11+num_subs)

	// The number of subs is not constant, therefore there is
	// a need for some smart assignment. The following array
	// sets the positions of thr first 5 subs, and then iterates
	// cyclicly. For example, if there are 2 subs allowed,
	// their positions will be GK (mandatory 1st !) and MF
	// If 7: GK, DF, MF, DF, FW, MF, DF
	//                              ^
	//                              cyclic repetition begins
	//
	sub_position := []string{"DFC", "MFC", "DFC", "FWC", "MFC"}

	// Iterates (cyclicly) over positions of subs,
	//
	sub_pos_iter := 0

	if len(roster.Players) < 11+num_subs {
		return internal.Teamsheet{}, fmt.Errorf("not enough players in roster")
	}

	ts := internal.Teamsheet{
		Name:   roster.TeamName,
		Code:   roster.TeamCode,
		Tactic: tactic,
	}

	var dfs, mfs, fws int
	var tact string

	if err := parse_formation(tactic, &dfs, &mfs, &fws, &tact); err != nil {
		return ts, err
	}

	// Calculate indices of the last defender and the last midfielder
	//
	last_df := dfs
	last_mf := dfs + mfs

	// Pick the players
	//
	// First, the best shot stopper is picked as a GK, then
	// others are picker according to the schedule of sub_position
	// as described above
	//

	// This will keep us from picking the same players more than once
	//
	chosen_players := make(map[string]any)

	for i := 0; i <= 10; i++ {
		if i == 0 {
			t_player[i].Pos = "GK"
		} else if i >= 1 && i <= last_df {
			t_player[i].Pos = "DFC"
		} else if i > last_df && i <= last_mf {
			t_player[i].Pos = "MFC"
		} else if i > last_mf && i <= 10 {
			t_player[i].Pos = "FWC"
		}
	}

	// set the best GK for N.1 position
	//
	t_player[0].Name = choose_best_player(roster.Players, chosen_players, st_getter)
	chosen_players[t_player[0].Name] = struct{}{}

	// From now on, j is the index for players in the teamsheet
	//

	// Set the starting defenders
	//
	for j := 1; j <= last_df; j++ {
		t_player[j].Name = choose_best_player(roster.Players, chosen_players, tk_getter)
		chosen_players[t_player[j].Name] = struct{}{}
	}

	// Set the starting midfielders
	//
	for j := last_df + 1; j <= last_mf; j++ {
		t_player[j].Name = choose_best_player(roster.Players, chosen_players, ps_getter)
		chosen_players[t_player[j].Name] = struct{}{}
	}

	// Set the starting forwards
	//
	for j := last_mf + 1; j < 11; j++ {
		t_player[j].Name = choose_best_player(roster.Players, chosen_players, sh_getter)
		chosen_players[t_player[j].Name] = struct{}{}
	}

	// Set the substitute GK
	//
	t_player[11].Name = choose_best_player(roster.Players, chosen_players, st_getter)
	t_player[11].Pos = "GK"
	chosen_players[t_player[11].Name] = struct{}{}

	name_of_best := ""

	for j := 12; j < num_subs+11; j++ {
		// What position should the current sub be on ?
		//
		if sub_position[sub_pos_iter] == "DFC" {
			name_of_best = choose_best_player(roster.Players, chosen_players, tk_getter)
		} else if sub_position[sub_pos_iter] == "MFC" {
			name_of_best = choose_best_player(roster.Players, chosen_players, ps_getter)
		} else if sub_position[sub_pos_iter] == "FWC" {
			name_of_best = choose_best_player(roster.Players, chosen_players, sh_getter)
		} else {
			return ts, fmt.Errorf("position index out of bound")
		}

		t_player[j].Name = name_of_best
		t_player[j].Pos = sub_position[sub_pos_iter]
		chosen_players[t_player[j].Name] = struct{}{}
		sub_pos_iter = (sub_pos_iter + 1) % 5
	}

	ts.Field = t_player[0:11]
	ts.Bench = t_player[11:]
	ts.PenaltyTaker = t_player[last_mf].Name

	return ts, nil
}

type skillCalculator func(internal.Player) int

func st_getter(player internal.Player) int {
	return player.St * player.Fitness / 100
}

func tk_getter(player internal.Player) int {
	return player.Tk * player.Fitness / 100
}

func ps_getter(player internal.Player) int {
	return player.Ps * player.Fitness / 100
}

func sh_getter(player internal.Player) int {
	return player.Sh * player.Fitness / 100
}

// / Gets the best player on some position from an array of roster players.
// /
// / players 		- the array of players
// / chosen_players 	- a set of already chosen players (those won't be chosen again)
// / skill 			- pointer to a function receiving a player and returning the skill by
// / 				  which "best" is judged.
// /
// / Returns the chosen player's name. Note: chosen_players is not modified !
// /
func choose_best_player(players []internal.Player, chosen_players map[string]any, skill skillCalculator) string {
	best_skill := -1
	name_of_best := ""

	for _, player := range players {
		if _, ok := chosen_players[player.Name]; ok {
			// The player has already been chosen
			continue
		}

		if player.Injury > 0 || player.Suspension > 0 {
			// The player cannot be selected
			continue
		}

		if player_skill := skill(player); player_skill > best_skill {
			best_skill = player_skill
			name_of_best = player.Name
		}
	}

	if name_of_best == "" {
		panic(fmt.Errorf("could not find best player"))
	}

	return name_of_best
}

// Parses the formation line, finds out how many defenders,
// midfielders and forwards to pick, and the tactic to use,
// performs error checking
//
// For example: 442N means 4 DFs, 4 MFs, 2 FWs, playing N
func parse_formation(formation string, dfs, mfs, fws *int, tactic *string) error {
	if len(formation) != 4 {
		return fmt.Errorf("the formation string must be exactly 4 characters long; for example: 442N")
	}

	// Random formation ?
	//
	if formation == "rand" {
		// between 3 and 5
		*dfs = 3 + rnd.Int()%3

		// if there are 5 dfs, max of 4 mfs
		if *dfs == 5 {
			*mfs = 1 + rnd.Int()%4
		} else { // 5 mfs is also possible
			*mfs = 1 + rnd.Int()%5
		}

		*fws = 10 - *dfs - *mfs

		*tactic = formation[3:]

		return nil
	}

	*dfs, _ = strconv.Atoi(formation[0:1])
	*mfs, _ = strconv.Atoi(formation[1:2])
	*fws, _ = strconv.Atoi(formation[2:3])

	*tactic = formation[3:]

	verify_position_range(*dfs)
	verify_position_range(*mfs)
	verify_position_range(*fws)

	if *dfs+*mfs+*fws != 10 {
		return fmt.Errorf("the number of players on all positions added together must be 10; for example: 442N")
	}

	return nil
}

func verify_position_range(n int) {
	if n < 1 || n > 8 {
		panic(fmt.Errorf("The number of players on each position must be between 1 and 8; For example: 442N"))
	}
}
