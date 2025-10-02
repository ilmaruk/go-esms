package esms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const (
	YELLOW = 1
	RED    = 2
)

var (
	numPlayers int

	theConfig                = &Config{}
	team_stats_total_enabled = false

	tact_manager = &TacticManager{}

	teams        [2]Team // two teams
	injuredInd   [2]int  // indicators of injured players for both teams
	yellowCarded [2]int  // indicators of yellow carded players for both teams
	redCarded    [2]int  // indicators of red carded players for both teams

	injured_ind [2]int // indices of injured players for both teams

	homeBonus float64

	teamStatsTotal [2][][3]float64 // [team][minute_index][tk,ps,sh]

	rnd *rand.Rand // random number generator

	substitutions int
	injuries      int
	fouls         int

	comm io.Writer = os.Stdout
)

type DidWhat int

const (
	DID_SHOT DidWhat = iota
	DID_FOUL
	DID_TACKLE
	DID_ASSIST
)

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func Play() {
	var (
		homeTeamsheet Teamsheet
		awayTeamsheet Teamsheet
	)

	dataDir := "../../data"

	// Home Teamsheet
	if err := loadTeamsheet(filepath.Join(dataDir, "ss1_sht.json"), &homeTeamsheet); err != nil {
		panic(err)
	}
	teams[0].Name = homeTeamsheet.Name

	// Away Teamsheet
	if err := loadTeamsheet(filepath.Join(dataDir, "ss2_sht.json"), &awayTeamsheet); err != nil {
		panic(err)
	}
	teams[1].Name = awayTeamsheet.Name

	var (
		homeRoster Roster
		awayRoster Roster
	)

	// Home Teamsheet
	if err := loadRoster(filepath.Join(dataDir, "ss1.json"), &homeRoster); err != nil {
		panic(err)
	}
	teams[0].RosterPlayer = homeRoster.Players

	// Away Teamsheet
	if err := loadRoster(filepath.Join(dataDir, "ss2.json"), &awayRoster); err != nil {
		panic(err)
	}
	teams[1].RosterPlayer = awayRoster.Players

	numSubs := theConfig.getIntConfig("NUM_SUBS", 5)
	numPlayers = 11 + numSubs

	team_stats_total_enabled = theConfig.getIntConfig("TEAM_STATS_TOTAL", 0) == 1 // TODO: default should be 0

	init_teams_data([2]Teamsheet{homeTeamsheet, awayTeamsheet})

	//--------------------------------------------
	//---------- The game running loop -----------
	//--------------------------------------------
	//
	// The timing logic is as follows:
	//
	// The game is divided to two structurally identical
	// halves. The difference between the halves is their
	// start times.
	//
	// For each half, an injury time is added. This time
	// goes into the minute counter, but not into the
	// formal_minute counter (that is needed for reports)
	//

	const half_length = 45

	// For each half
	//
	for half_start := 1; half_start < 2*half_length; half_start += half_length {
		half := 1
		if half_start != 1 {
			half = 2
		}
		last_minute_of_half := half_start + half_length - 1
		in_inj_time := false

		// Play the game minutes of this half
		//
		// last_minute_of_half will be increased by inj_time_length in
		// the end of the half
		//
		formal_minute := half_start
		for minute := formal_minute; minute <= last_minute_of_half; minute++ {
			cleanInjCardIndicators()
			recalculate_teams_data()

			// For each team
			//
			for j := 0; j <= 1; j++ {
				// Calculate different events
				//
				ifShot(j, minute)
				ifFoul(j)
				randomInjury(j)

				// notJ := 1 - j
				// score_diff := teams[j].Score - teams[notJ].Score
				// check_conditionals(j);
			}

			// fixme ?
			if team_stats_total_enabled {
				if minute == 1 || minute%10 == 0 {
					add_team_stats_total(minute)
				}
			}

			if !in_inj_time {
				formal_minute++

				updatePlayersMinuteCount()
			}

			if minute == last_minute_of_half && !in_inj_time {
				in_inj_time = true

				// shouldn't have been increased, but we only know about
				// this now
				formal_minute--

				inj_time_length := how_much_inj_time()
				last_minute_of_half += inj_time_length

				// char buf[2000];
				// snprintf(buf, 2000 * sizeof(char), "%d", inj_time_length);
				// fprintf(comm, "\n%s\n", the_commentary().rand_comment("COMM_INJURYTIME", buf).c_str());
			}
		}

		in_inj_time = false

		if half == 1 {
			// fprintf(comm, "\n%s\n", the_commentary().rand_comment("COMM_HALFTIME").c_str())
			fmt.Fprintf(comm, "\nCOMM_HALFTIME\n")
		} else if half == 2 {
			// fprintf(comm, "\n%s\n", the_commentary().rand_comment("COMM_FULLTIME").c_str())
			fmt.Fprintf(comm, "\nCOMM_FULLTIME\n")
		}
	}

	calcAbility()

	fmt.Fprintln(comm, "Final score:", teams[0].Name, teams[0].Score, "-", teams[1].Score, teams[1].Name)
}

// Temporary, for testing purposes
func prettyPrintTeam(t Team) {
	b, _ := json.Marshal(t)
	fmt.Println(string(b))

	var out bytes.Buffer
	if err := json.Indent(&out, b, "", "  "); err != nil {
		fmt.Fprintln(os.Stderr, "invalid JSON:", err)
		os.Exit(1)
	}

	out.WriteByte('\n')
	if _, err := out.WriteTo(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "write error:", err)
		os.Exit(1)
	}
}

// / Calculates how much injury time to add.
// /
// / Takes into account substitutions, injuries and fouls (by both teams)
// /
func how_much_inj_time() int {
	// Each time this function is called, it subtracts the last
	// totals it had, because the stats accumulate and don't
	// annulize between halves
	//
	substitutions := teams[0].Substitutions + teams[1].Substitutions - substitutions
	injuries := teams[0].Injuries + teams[1].Injuries - injuries
	fouls := teams[0].FinalFouls + teams[1].FinalFouls - fouls

	calc := math.Ceil(float64(substitutions)*0.5 + float64(injuries)*0.5 + float64(fouls)*0.5)

	return int(calc)
}

func add_team_stats_total(minute int) {
	/* Define the correct index for the array */
	var index int
	if minute == 1 {
		index = 0
	} else {
		index = minute / 10
	}

	for i := 0; i <= 1; i++ {
		teamStatsTotal[i][index][0] = teams[i].TeamTackling
		teamStatsTotal[i][index][1] = teams[i].TeamPassing
		teamStatsTotal[i][index][2] = teams[i].TeamShooting
	}
}

func init_teams_data(teamsheet [2]Teamsheet) {
	for l := 0; l <= 1; l++ {
		teams[l].Tactic = teamsheet[l].Tactic

		if !tact_manager.tactic_exists(teams[l].Tactic) {
			panic(fmt.Errorf("Invalid tactic %s in %s's teamsheet", teams[l].Tactic, teams[l].Name))
		}

		for i := 0; i < numPlayers; i++ {
			var full_pos string

			/* Read players's position and name */
			if i < 11 {
				teams[l].Players[i].Name = teamsheet[l].Field[i].Name
				full_pos = teamsheet[l].Field[i].Pos
			} else {
				teams[l].Players[i].Name = teamsheet[l].Bench[i-11].Name
				full_pos = teamsheet[l].Bench[i-11].Pos
			}

			// For GKs, just copy the position as is
			//
			if full_pos == "GK" {
				teams[l].Players[i].Pos = "GK"
			} else {
				if !is_legal_position(full_pos) {
					panic(fmt.Errorf("illegal position %s of %s in %s's teamsheet", full_pos, teams[l].Players[i].Name, teams[l].Name))
				}

				teams[l].Players[i].Pos = fullpos2position(full_pos)
				teams[l].Players[i].Side = fullpos2side(full_pos)
			}

			/* The first specified player must be a GK */
			if i == 0 && teams[l].Players[i].Pos != "GK" {
				panic(fmt.Errorf("the first player in %s's teamsheet must be a GK", teams[l].Name))
			}

			if teams[l].Players[i].Pos == "PK:" {
				panic(fmt.Errorf("pk: where player %d was expected (%s)", i, teams[l].Name))
			}

			found := 0

			// Search for this player in the roster, and when found assign his info
			// to the player structure.
			//
			for _, player := range teams[l].RosterPlayer {
				if teams[l].Players[i].Name != player.Name {
					continue
				}

				found = 1

				// Check if the player is available for the game
				//
				if player.Injury > 0 {
					panic(fmt.Errorf("Player %s (%s) is injured", player.Name, teams[l].Name))
				}

				if player.Suspension > 0 {
					panic(fmt.Errorf("Player %s (%s) is suspended", player.Name, teams[l].Name))
				}

				teams[l].Players[i].pref_side = player.PrefSide

				teams[l].Players[i].likes_left = false
				teams[l].Players[i].likes_right = false
				teams[l].Players[i].likes_center = false

				if teams[l].Players[i].pref_side == "L" {
					teams[l].Players[i].likes_left = true
				}

				if teams[l].Players[i].pref_side == "R" {
					teams[l].Players[i].likes_right = true
				}

				if teams[l].Players[i].pref_side == "C" {
					teams[l].Players[i].likes_center = true
				}

				teams[l].Players[i].st = player.St
				teams[l].Players[i].tk = player.Tk
				teams[l].Players[i].ps = player.Ps
				teams[l].Players[i].sh = player.Sh
				teams[l].Players[i].stamina = player.Stamina

				// Each player has a nominal_fatigue_per_minute rating that's
				// calculated once, based on his stamina.
				//
				// I'd like the average rating be 0.031 - so that an average player
				// (stamina = 50) will lose 30 fitness points during a full game.
				//
				// The range is approximately 50 - 10 points, and the stamina range
				// is 1-99. So, first the ratio is normalized and then subtracted
				// from the average 0.031 (which, times 90 minutes, is 0.279).
				// The formula for each player is:
				//
				// fatigue            stamina - 50
				// ------- = 0.0031 - ------------  * 0.0022
				//  minute                 50
				//
				//
				// This gives (approximately) 30 lost fitness points for average players,
				// 50 for the worse stamina and 10 for the best stamina.
				//
				// A small random factor is added each minute, so the exact numbers are
				// not deterministic.
				//
				normalized_stamina_ratio := float64(teams[l].Players[i].stamina-50) / 50.0
				teams[l].Players[i].nominal_fatigue_per_minute = 0.0031 - normalized_stamina_ratio*0.0022

				teams[l].Players[i].ag = player.Ag
				teams[l].Players[i].fatigue = float64(player.Fitness) / 100.0

				break
			}

			if found == 0 {
				panic(fmt.Errorf("Player %s (%s) doesn't exist in the roster file", teams[l].Players[i].Name, teams[l].Name))
			}
		}

		// There's an optional "PK: <Name>" line.
		// If it exists, the <Name> must be listed in the teamsheet.
		var i int
		for i = numPlayers - 1; i >= 0; i-- {
			if teamsheet[l].PK == teams[l].Players[i].Name {
				teams[l].PenaltyTaker = i
				break
			}
		}

		if i < 0 {
			panic(fmt.Errorf("error in penalty kick taker of %s, player %s not listed", teams[l].Name, teamsheet[l].PK))
		}
	}

	ensure_no_duplicate_names()

	// Temporarily disable
	// read_conditionals(teamsheet);

	// Set active flags
	for j := 0; j <= 1; j++ {
		teams[j].Substitutions = 0
		teams[j].Injuries = 0

		for i := 0; i < numPlayers; i++ {
			if i < 11 {
				teams[j].Players[i].Active = 1
			} else {
				teams[j].Players[i].Active = 2
			}
		}
	}

	/* In the beginning, player n.1 is always the GK */
	teams[0].CurrentGK = 1
	teams[1].CurrentGK = 1

	/* Data initialization */
	for j := 0; j <= 1; j++ {
		teams[j].Score = 0
		teams[j].FinalShotsOn = 0
		teams[j].FinalShotsOff = 0
		teams[j].FinalFouls = 0
		teams[j].TeamTackling = 0
		teams[j].TeamPassing = 0
		teams[j].TeamShooting = 0

		for i := 0; i < numPlayers; i++ {
			teams[j].Players[i].tk_contrib = 0
			teams[j].Players[i].ps_contrib = 0
			teams[j].Players[i].sh_contrib = 0

			teams[j].Players[i].yellowcards = 0
			teams[j].Players[i].redcards = 0
			teams[j].Players[i].injured = 0
			teams[j].Players[i].tk_ab = 0
			teams[j].Players[i].ps_ab = 0
			teams[j].Players[i].sh_ab = 0
			teams[j].Players[i].st_ab = 0

			// final stats initialization
			teams[j].Players[i].Minutes = 0
			teams[j].Players[i].shots = 0
			teams[j].Players[i].goals = 0
			teams[j].Players[i].saves = 0
			teams[j].Players[i].assists = 0
			teams[j].Players[i].tackles = 0
			teams[j].Players[i].keypasses = 0
			teams[j].Players[i].fouls = 0
			teams[j].Players[i].redcards = 0
			teams[j].Players[i].yellowcards = 0
			teams[j].Players[i].conceded = 0
			teams[j].Players[i].shots_on = 0
			teams[j].Players[i].shots_off = 0
		}
	}
}

// / Goes over both teams and checks that there are no duplicate player
// / names. If there are, exits with an error.
// /
func ensure_no_duplicate_names() {
	for j := 0; j <= 1; j++ {
		for i := 0; i <= numPlayers; i++ {
			for k := 0; k <= numPlayers; k++ {
				if k != i && teams[j].Players[i].Name == teams[j].Players[k].Name {
					panic(fmt.Errorf("Player %s (%s) is named twice in the team sheet", teams[j].Players[i].Name, teams[j].Name))
				}
			}
		}
	}
}

func change_tactic(a int, newtct string) {
	if newtct != teams[a].Tactic {
		teams[a].Tactic = newtct

		// fputs(the_commentary().rand_comment("CHANGETACTIC",
		//                                     minute_str().c_str(),
		//                                     teams[a].name, teams[a].name,
		//                                     teams[a].tactic)
		//           .c_str(),
		//       comm);
	}
}

/* This function controls the random injuries occurance. */
/* The CHANCE of a player to get injured depends on a    */
/* constant factor + total aggression of the rival team. */
/* The function will find who was injured and substitute */
/* him for player on his position.                       */
func randomInjury(a int) {
	var injured int
	var found int

	notA := 1 - a
	if randomp(int((1500+teams[notA].Aggression)/50)) == 1 { /* If someone got injured */
		teams[a].Injuries++

		for { /* The inj_player can't be n.0 and must be playing */
			injured = myRandom(numPlayers + 1)
			if injured != 0 && teams[a].Players[injured].Active == 1 {
				break
			}
		} // while (injured == 0 || teams[a].Players[injured].active != 1);

		// fprintf(comm, "%s",
		//         the_commentary().rand_comment("INJURY", minute_str().c_str(), teams[a].name,
		//                                       teams[a].Players[injured].name.c_str())
		//             .c_str());

		// report_event *an_event = new report_event_injury(teams[a].Players[injured].name,
		//                                                  teams[a].name, formal_minute_str().c_str());
		// report_vec.push_back(an_event);

		injured_ind[a] = injured

		/* Only 3 substitutions are allowed per team per game */
		if teams[a].Substitutions >= 5 { /* No substitutions left */
			teams[a].Players[injured].Active = 0
			// fprintf(comm, "%s", the_commentary().rand_comment("NOSUBSLEFT").c_str());

			if teams[a].Players[injured].Pos == "GK" {
				n := 11

				for teams[a].Players[n].Active != 1 { /* Sub him for another player */
					n--
				}

				changePosition(a, n, string("GK"))
				teams[a].CurrentGK = n
			}
		} else {
			b := 12

			for found == 0 && b < numPlayers { /* Look for subs on the same position */
				if teams[a].Players[injured].Pos == teams[a].Players[b].Pos && teams[a].Players[b].Active == 2 {
					substitutePlayer(a, injured, b, posAndSide2fullpos(teams[a].Players[injured].Pos, teams[a].Players[injured].Side))

					if injured == teams[a].CurrentGK {
						teams[a].CurrentGK = b
					}

					found = 1
				} else {
					b++
				}
			}

			if found == 0 { /* If there are no subs on his position */
				/* Then, sub him for any other player on the bench who is not a   */
				/* goalkeeper. If a GK will be injured, he will be subbed for the */
				/* GK on the bench by the previous loop, if there won't be any    */
				/* GK on the bench, he will be subbed for another player          */
				b = 12

				for found == 0 && b < numPlayers {
					if teams[a].Players[b].Pos != "GK" && teams[a].Players[b].Active == 2 {
						substitutePlayer(a, injured, b, posAndSide2fullpos(teams[a].Players[injured].Pos, teams[a].Players[injured].Side))
						found = 1

						if injured == teams[a].CurrentGK {
							teams[a].CurrentGK = b
						}
					} else {
						b++
					}
				} // while (!found && b <= num_players)
			} // if (!found)
		} // if (teams[a].substitutions >= 3)

		teams[a].Players[injured].injured = 1
		teams[a].Players[injured].Active = 0

	} // if (randomp((1500 + teams[notA].aggression)/50))
}

func calc_shotprob(a int) {
	// Note: 1.0 is added to tackling, to avoid singularity when the
	// team tackling is 0
	//
	notA := 1 - a
	teams[a].ShotProb = 1.8 * (teams[a].Aggression/50.0 + 800.0*math.Pow(((1.0/3.0*teams[a].TeamShooting+2.0/3.0*teams[a].TeamPassing)/(teams[notA].TeamTackling+1.0)), 2))

	// If it is the home team, add home bonus
	//
	if a == 0 {
		teams[a].ShotProb += homeBonus
	}
}

// Calculate the contributions of player b of team a
func calc_player_contributions(a, b int) {
	notA := 1 - a
	if teams[a].Players[b].Active == 1 && teams[a].CurrentGK != b {
		tk_mult := tact_manager.get_mult(teams[a].Tactic, teams[notA].Tactic, teams[a].Players[b].Pos, "TK")
		ps_mult := tact_manager.get_mult(teams[a].Tactic, teams[notA].Tactic, teams[a].Players[b].Pos, "PS")
		sh_mult := tact_manager.get_mult(teams[a].Tactic, teams[notA].Tactic, teams[a].Players[b].Pos, "SH")

		var side_factor float64

		if (teams[a].Players[b].Side == "R" && teams[a].Players[b].likes_right) ||
			(teams[a].Players[b].Side == "L" && teams[a].Players[b].likes_left) ||
			(teams[a].Players[b].Side == "C" && teams[a].Players[b].likes_center) {
			side_factor = 1.0
		} else {
			side_factor = 0.75
		}

		teams[a].Players[b].tk_contrib = tk_mult * side_factor * float64(teams[a].Players[b].tk) * teams[a].Players[b].fatigue
		teams[a].Players[b].ps_contrib = ps_mult * side_factor * float64(teams[a].Players[b].ps) * teams[a].Players[b].fatigue
		teams[a].Players[b].sh_contrib = sh_mult * side_factor * float64(teams[a].Players[b].sh) * teams[a].Players[b].fatigue
	} else { // The contributions of an inactive player or of a GK are 0
		teams[a].Players[b].tk_contrib = 0
		teams[a].Players[b].ps_contrib = 0
		teams[a].Players[b].sh_contrib = 0
	}
}

// Adjusts players' total contributions, taking into account the
// side balance on each position
func adjust_contrib_with_side_balance(a int) {
	// The side balance:
	// For each position (w/o side), keep a vector of 3 elements
	// to specify the number of players playing R [0], L [1], C [2] on this position
	//
	balance := make(map[string][]int)

	// Init the side balance for all positions
	//
	positions := tact_manager.get_positions_names()
	for _, pos := range positions {
		balance[pos] = []int{0, 0, 0}
	}

	// Go over the team's players and record on what side they play,
	// updating the side balance
	//
	for b := 1; b < numPlayers; b++ {
		if teams[a].Players[b].Active == 1 && teams[a].Players[b].Pos != "GK" {
			if teams[a].Players[b].Side == "R" {
				balance[teams[a].Players[b].Pos][0]++
			} else if teams[a].Players[b].Side == "L" {
				balance[teams[a].Players[b].Pos][1]++
			} else if teams[a].Players[b].Side == "C" {
				balance[teams[a].Players[b].Pos][2]++
			} else {
				panic(fmt.Errorf("internal error"))
			}
		}
	}

	// For all positions, check if the side balance is equal for R and L
	// If it isn't, penalize the contributions of the players on those positions
	//
	// Additionally, penalize teams who play with more than 3 C players on
	// some position without R and L
	//
	for _, pos := range positions {
		on_pos_right := balance[pos][0]
		on_pos_left := balance[pos][1]
		on_pos_center := balance[pos][2]

		var taxed_multiplier float64 = 1.0

		if on_pos_left != on_pos_right {
			tax_ratio := 0.25 * math.Abs(float64(on_pos_right-on_pos_left)) / float64(on_pos_right+on_pos_left)
			taxed_multiplier = 1 - tax_ratio
		} else if on_pos_left == 0 && on_pos_right == 0 && on_pos_center > 3 {
			taxed_multiplier = 0.87
		}

		if taxed_multiplier != 1 {
			for b := 1; b < numPlayers; b++ {
				if teams[a].Players[b].Active == 1 && teams[a].Players[b].Pos == pos {
					teams[a].Players[b].tk_contrib *= taxed_multiplier
					teams[a].Players[b].ps_contrib *= taxed_multiplier
					teams[a].Players[b].sh_contrib *= taxed_multiplier
				}
			}
		}
	}
}

// This function is called by the game running loop in the
// beginning of each minute of the game.
// It recalculates player contributions, aggression, fatigue,
// team total contributions and shotprob.
func recalculate_teams_data() {
	for a := 0; a <= 1; a++ {
		teams[a].TeamTackling = 0
		teams[a].TeamPassing = 0
		teams[a].TeamShooting = 0
		calc_aggression(a)

		for b := 1; b < numPlayers; b++ {
			if teams[a].Players[b].Active == 1 {
				fatigue_deduction := teams[a].Players[b].nominal_fatigue_per_minute
				mrnd := myRandom(100)
				fatigue_deduction += float64(mrnd-50) / 50.0 * 0.003

				teams[a].Players[b].fatigue -= fatigue_deduction

				if teams[a].Players[b].fatigue < 0.10 {
					teams[a].Players[b].fatigue = 0.10
				}
			}
		}

		for b := 1; b < numPlayers; b++ {
			calc_player_contributions(a, b)
		}

		adjust_contrib_with_side_balance(a)
		calc_team_contributions_total(a)
	}

	for a := 0; a <= 1; a++ {
		calc_shotprob(a)
	}
}

func calc_team_contributions_total(a int) {
	for b := 2; b <= numPlayers; b++ {
		if teams[a].Players[b].Active == 1 {
			teams[a].TeamTackling += teams[a].Players[b].tk_contrib
			teams[a].TeamPassing += teams[a].Players[b].ps_contrib
			teams[a].TeamShooting += teams[a].Players[b].sh_contrib
		}
	}
}

// This function sets the aggression of all inactive players to 0
// and then adds up all aggressions in the team total aggression
func calc_aggression(a int) {
	teams[a].Aggression = 0

	for i := 0; i < numPlayers; i++ {
		if teams[a].Players[i].Active != 1 {
			teams[a].Players[i].ag = 0
		}

		teams[a].Aggression += teams[a].Players[i].ag
	}
}

// Called on each minute to handle a scoring chance of team
// a for this minute.
func ifShot(a, minute int) {
	var shooter int
	var assister int
	var tackler int
	var chance_tackled int
	chance_assisted := 0

	// Did a scoring chance occur ?
	//
	if randomp(int(teams[a].ShotProb)) == 1 {
		// There's a 0.75 probability that a chance was assisted, and
		// 0.25 that it's a solo
		//
		if randomp(7500) == 1 {
			assister = whoDidIt(a, DID_ASSIST)
			chance_assisted = 1

			shooter = whoGotAssist(a, assister)

			// fprintf(comm, "%s", the_commentary().rand_comment("ASSISTEDCHANCE", minute_str().c_str(), teams[a].name, teams[a].Players[assister].name.c_str(), teams[a].Players[shooter].name.c_str()).c_str());
			fmt.Fprintln(comm, "ASSISTEDCHANCE", minute, teams[a].Name, teams[a].Players[assister].Name, teams[a].Players[shooter].Name)
			teams[a].Players[assister].keypasses++
		} else {
			shooter = whoDidIt(a, DID_SHOT)

			chance_assisted = 0
			assister = 0

			// fprintf(comm, "%s", the_commentary().rand_comment("CHANCE", minute_str().c_str(), teams[a].name, teams[a].Players[shooter].name.c_str()).c_str());
			fmt.Fprintln(comm, "CHANCE", minute, teams[a].Name, teams[a].Players[shooter].Name)
		}

		notA := 1 - a
		chance_tackled = int(4000.0 * ((teams[notA].TeamTackling * 3.0) / (teams[a].TeamPassing*2.0 + teams[a].TeamShooting)))

		/* If the chance was tackled */
		if randomp(chance_tackled) == 1 {
			tackler = whoDidIt(notA, DID_TACKLE)
			teams[notA].Players[tackler].tackles++

			// fprintf(comm, "%s", the_commentary().rand_comment("TACKLE", teams[notA].Players[tackler].name.c_str()).c_str());
			fmt.Fprintln(comm, "\tTACKLE", teams[notA].Players[tackler].Name)
		} else { /* Chance was not tackled, it will be a shot on goal */
			// fprintf(comm, "%s", the_commentary().rand_comment("SHOT", teams[a].Players[shooter].name.c_str()).c_str());
			fmt.Fprintln(comm, "\tSHOT", teams[a].Players[shooter].Name)
			teams[a].Players[shooter].shots++

			if ifOnTarget(a, shooter) == 1 {
				teams[a].FinalShotsOn++
				teams[a].Players[shooter].shots_on++

				if ifGoal(a, shooter) == 1 {
					// fprintf(comm, "%s", the_commentary().rand_comment("GOAL").c_str());
					fmt.Fprintln(comm, "\t\tGOAL")

					if isGoalCancelled() == 0 {
						teams[a].Score++

						// If the assister was the shooter, there was no
						// assist, but a simple goal.
						//
						if chance_assisted == 1 && (assister != shooter) {
							teams[a].Players[assister].assists++ /* For final stats */
						}

						teams[a].Players[shooter].goals++
						teams[notA].Players[teams[notA].CurrentGK].conceded++

						// fprintf(comm, "\n          ...  %s %d-%d %s ...",
						//         teams[0].name,
						//         teams[0].score,
						//         teams[1].score,
						//         teams[1].name);

						// report_event *an_event = new report_event_goal(teams[a].Players[shooter].name.c_str(),
						//                                                teams[a].name, formal_minute_str().c_str());

						// report_vec.push_back(an_event);
					}
				} else {
					// fprintf(comm, "%s", the_commentary().rand_comment("SAVE", teams[notA].Players[teams[notA].current_gk].name.c_str()).c_str());
					fmt.Fprintln(comm, "\t\tSAVE", teams[notA].Players[teams[notA].CurrentGK].Name)

					teams[notA].Players[teams[notA].CurrentGK].saves++
				}
			} else {
				teams[a].Players[shooter].shots_off++
				// fprintf(comm, "%s", the_commentary().rand_comment("OFFTARGET").c_str())
				fmt.Fprintln(comm, "\tOFFTARGET")
				teams[a].FinalShotsOff++
			}
		}
	}
}

// When a chance was generated for the team and assisted by the
// assister, who got the assist ?
//
// This is almost like who_did_it, but it also takes
// into account the side of the assister - a player on his side
// has a higher chance to get the assist.
//
// How it's done: if the side of the shooter (picked by who_did_it)
// is different from the side of the asssiter, who_did_it is run
// once again - but this happens only once. This increases the
// chance of the player on the same side to be picked, but leaves
// a possibility for other sides as well.
func whoGotAssist(a, assister int) int {
	shooter := assister

	// Shooter and assister must be different, so re-run each time the same
	// one is generated
	//
	for shooter == assister {
		shooter = whoDidIt(a, DID_SHOT)

		// if the side is different, re-run once
		//
		if teams[a].Players[shooter].Side != teams[a].Players[assister].Side {
			shooter = whoDidIt(a, DID_SHOT)
		}
	}

	return shooter
}

/* Whether the shot is on target. */
func ifOnTarget(a, b int) int {
	if randomp(int(5800.0*teams[a].Players[b].fatigue)) == 1 {
		return 1
	} else {
		return 0
	}
}

// Given a shot on target (team a shot on team b's goal),
// was it a goal ?
func ifGoal(a, b int) int {
	// Factors taken into account:
	// The shooter's Sh and fatigue against the GK's St
	//
	// The "median" is 0.35
	// Lower and upper bounds are 0.1 and 0.9 respectively
	//
	notA := 1 - a
	var temp float64 = float64(teams[a].Players[b].sh*int(teams[a].Players[b].fatigue)*200 - teams[notA].Players[teams[notA].CurrentGK].st*200 + 3500)

	if temp > 9000 {
		temp = 9000
	}
	if temp < 1000 {
		temp = 1000
	}

	if randomp(int(temp)) == 1 {
		return 1
	} else {
		return 0
	}
}

func isGoalCancelled() int {
	if randomp(500) == 1 {
		// fprintf(comm, "%s", the_commentary().rand_comment("GOALCANCELLED").c_str());
		return 1
	}

	return 0
}

// Given a team and an event (eg. SHOT)
// picks one player at (weighted) random
// that performed this event.
//
// For example, for SHOT, pick a player
// at weighted random according to his
// shooting skill
func whoDidIt(a int, event DidWhat) int {
	k := 0
	var total float64 = 0
	var weight float64 = 0
	ar := make([]float64, numPlayers)

	// Employs the weighted random algorithm
	// A player's chance to DO_IT is his
	// contribution relative to the team's total
	// contribution
	//

	for k = 0; k < numPlayers; k++ {
		switch event {
		case DID_SHOT:
			weight += teams[a].Players[k].sh_contrib * 100.0
			total = teams[a].TeamShooting * 100.0
		case DID_FOUL:
			weight += teams[a].Players[k].ag
			total = teams[a].Aggression
		case DID_TACKLE:
			weight += teams[a].Players[k].tk_contrib * 100.0
			total = teams[a].TeamTackling * 100.0
		case DID_ASSIST:
			weight += teams[a].Players[k].ps_contrib * 100.0
			total = teams[a].TeamPassing * 100.0
		default:
			panic(fmt.Errorf("internal error"))
			// cout << "Internal error, " << __FILE__ << ", line " << __LINE__ << endl;
			// MY_EXIT(1);
		}

		ar[k] = weight
	}

	rand_value := float64(myRandom(int(total)))

	for k = 2; ar[k] <= rand_value; k++ {
		if k == numPlayers {
			panic(fmt.Errorf("internal error"))
			// cout << "Internal error, " << __FILE__ << ", line " << __LINE__ << endl;
			// MY_EXIT(1);
		}
	}

	// delete[] ar;

	return k
}

// ifFoul handles fouls (called on each minute with for each team)
func ifFoul(a int) {
	var fouler int

	if randomp(int(teams[a].Aggression*.75)) == 1 {
		fouler = whoDidIt(a, DID_FOUL)
		// fprintf(comm, "%s", the_commentary().rand_comment("FOUL", minute_str().c_str(), teams[a].name, teams[a].Players[fouler].name.c_str()).c_str());

		teams[a].FinalFouls++ /* For final stats */
		teams[a].Players[fouler].fouls++

		/* The chance of the foul to result in a yellow or red card */
		if randomp(6000) == 1 {
			bookings(a, fouler, YELLOW)
		} else if randomp(400) == 1 {
			bookings(a, fouler, RED)
		} else {
			// fprintf(comm, "%s", the_commentary().rand_comment("WARNED").c_str());
		}

		notA := 1 - a

		/* Condition for a penalty to occur (if GK fouled, or random) */
		if fouler == teams[a].CurrentGK || randomp(500) == 1 {
			// If the nominated PK taker isn't active, choose the
			// best shooter to take the PK
			//
			if teams[notA].Players[teams[notA].PenaltyTaker].Active != 1 || teams[notA].PenaltyTaker == -1 {
				var max float64 = -1
				max_index := 1

				for i := 0; i < numPlayers; i++ {
					if teams[notA].Players[i].Active == 1 && float64(teams[notA].Players[i].sh)*teams[notA].Players[i].fatigue > max {
						max = float64(teams[notA].Players[i].sh) * teams[notA].Players[i].fatigue
						max_index = i
					}
				}

				teams[notA].PenaltyTaker = max_index
			}

			// fprintf(comm, "%s", the_commentary().rand_comment("PENALTY", teams[notA].Players[teams[notA].penalty_taker].name.c_str()).c_str());

			/* If Penalty... Goal ? */
			if randomp(8000+teams[notA].Players[teams[notA].PenaltyTaker].sh*100-teams[a].Players[teams[a].CurrentGK].st*100) == 1 {
				// fprintf(comm, "%s", the_commentary().rand_comment("GOAL").c_str());
				teams[notA].Score++
				teams[notA].Players[teams[notA].PenaltyTaker].goals++
				teams[a].Players[teams[a].CurrentGK].conceded++
				// fprintf(comm, "\n          ...  %s %d-%d %s...", teams[0].name, teams[0].score,
				//         teams[1].score, teams[1].name);

				// report_event *an_event = new report_event_penalty(teams[notA].Players[teams[notA].penalty_taker].name,
				//                                                   teams[notA].name, formal_minute_str().c_str());
				// report_vec.push_back(an_event);
			} else { /* If the penalty taker didn't score */
				// Either it was saved, or it went off-target
				//
				if randomp(7500) == 1 {
					// fprintf(comm, "%s", the_commentary().rand_comment("SAVE", teams[a].Players[teams[a].current_gk].name.c_str()).c_str());
				} else { /* Or it went off-target */
					// fprintf(comm, "%s", the_commentary().rand_comment("OFFTARGET").c_str());
				}
			}
		}
	}
}

// bookings deals with yellow and red cards
func bookings(a, b, card_color int) {
	if card_color == YELLOW {
		// fprintf(comm, "%s", the_commentary().rand_comment("YELLOWCARD").c_str());
		teams[a].Players[b].yellowcards++

		// A second yellow card is equal to a red card
		//
		if teams[a].Players[b].yellowcards == 2 {
			// fprintf(comm, "%s", the_commentary().rand_comment("SECONDYELLOWCARD").c_str());
			sendOff(a, b)

			// report_event *an_event = new report_event_red_card(teams[a].Players[b].name.c_str(),
			//                                                    teams[a].name, formal_minute_str().c_str());
			// report_vec.push_back(an_event);

			redCarded[a] = b
		} else {
			yellowCarded[a] = b
		}
	} else if card_color == RED {
		// fprintf(comm, "%s", the_commentary().rand_comment("REDCARD").c_str());
		sendOff(a, b)

		// report_event *an_event = new report_event_red_card(teams[a].Players[b].name.c_str(),
		//                                                    teams[a].name, formal_minute_str().c_str());
		// report_vec.push_back(an_event);

		redCarded[a] = b
	}
}

// substitutePlayer substitutites player in for player out in team a, he'll play
// position newpos
func substitutePlayer(a, out, in int, newpos string) {
	max_substitutions := theConfig.getIntConfig("SUBSTITUTIONS", 5)

	if teams[a].Players[out].Active == 1 && teams[a].Players[in].Active == 2 && teams[a].Substitutions < max_substitutions {
		teams[a].Players[out].Active = 0
		teams[a].Players[in].Active = 1

		if newpos == "GK" {
			teams[a].Players[in].Pos = "GK"
		} else {
			teams[a].Players[in].Pos = fullpos2position(newpos)
			teams[a].Players[in].Side = fullpos2side(newpos)
		}

		if out == teams[a].CurrentGK {
			teams[a].CurrentGK = in
		}

		teams[a].Substitutions++

		// fputs(the_commentary().rand_comment("SUB", minute_str().c_str(), teams[a].name,
		//                                     teams[a].Players[in].name.c_str(),
		//                                     teams[a].Players[out].name.c_str(),
		//                                     newpos.c_str())
		//           .c_str(),
		//       comm);
	}
}

func changePosition(a, b int, newpos string) {
	// Can't reposition a GK or an inactive player
	if b != teams[a].CurrentGK && teams[a].Players[b].Active == 1 {
		// If he plays on this position anyway, don't change it
		if posAndSide2fullpos(teams[a].Players[b].Pos, teams[a].Players[b].Side) != newpos {
			// fputs(the_commentary().rand_comment("CHANGEPOSITION", minute_str().c_str(),
			//                                     teams[a].name,
			//                                     teams[a].Players[b].name.c_str(),
			//                                     newpos.c_str())
			//           .c_str(),
			//       comm);

			teams[a].Players[b].Pos = fullpos2position(newpos)
			teams[a].Players[b].Side = fullpos2side(newpos)
		}
	}
}

func sendOff(a, b int) {
	teams[a].Players[b].yellowcards = 0
	teams[a].Players[b].redcards++
	teams[a].Players[b].Active = 0

	if teams[a].CurrentGK == b { /* If a GK was sent off */
		i := 12
		found := false

		if teams[a].Substitutions < 3 {
			for !found && i <= numPlayers { /* Look for a keeper on the bench */
				/* If found a keeper */
				if teams[a].Players[i].Pos == "GK" && teams[a].Players[i].Active == 2 {
					n := 11

					found = true

					for teams[a].Players[n].Active != 1 { /* Sub him for another player */
						n--
					}
					substitutePlayer(a, n, i, "GK")
					teams[a].CurrentGK = i
				} else {
					found = false
					i++
				}
			}

			/*  If there was no keeper on the bench   */
			/*  Change the position of another player */
			/*  (who is on the field) to GK           */
			if !found {
				n := 11

				for teams[a].Players[n].Active != 1 {
					n--
				}

				changePosition(a, n, string("GK"))
				teams[a].CurrentGK = n
			}
		} else { /* If substitutions >= 3 */
			n := 11

			for teams[a].Players[n].Active != 1 {
				n--
			}
			changePosition(a, n, string("GK"))
			teams[a].CurrentGK = n
		}
	}
}

// calcAbility uses the constants contained in league.dat
// to calculate the ability change of each player.
func calcAbility() {
	// Initialization of ab bonuses
	//
	ab_goal := theConfig.getIntConfig("AB_GOAL", 0)
	ab_assist := theConfig.getIntConfig("AB_ASSIST", 0)
	ab_victory := theConfig.getIntConfig("AB_VICTORY_RANDOM", 0)
	ab_defeat := theConfig.getIntConfig("AB_DEFEAT_RANDOM", 0)
	ab_cleansheet := theConfig.getIntConfig("AB_CLEAN_SHEET", 0)
	ab_ktk := theConfig.getIntConfig("AB_KTK", 0)
	ab_kps := theConfig.getIntConfig("AB_KPS", 0)
	ab_sht_on := theConfig.getIntConfig("AB_SHT_ON", 0)
	ab_sht_off := theConfig.getIntConfig("AB_SHT_OFF", 0)
	ab_sav := theConfig.getIntConfig("AB_SAV", 0)
	ab_concede := theConfig.getIntConfig("AB_CONCDE", 0)
	ab_yellow := theConfig.getIntConfig("AB_YELLOW", 0)
	ab_red := theConfig.getIntConfig("AB_RED", 0)

	for j := 0; j <= 1; j++ {
		// Add simple bonuses
		//
		for i := 1; i < numPlayers; i++ {
			teams[j].Players[i].sh_ab += ab_goal * teams[j].Players[i].goals
			teams[j].Players[i].ps_ab += ab_assist * teams[j].Players[i].assists
			teams[j].Players[i].tk_ab += ab_ktk * teams[j].Players[i].tackles
			teams[j].Players[i].ps_ab += ab_kps * teams[j].Players[i].keypasses
			teams[j].Players[i].sh_ab += ab_sht_on * teams[j].Players[i].shots_on
			teams[j].Players[i].sh_ab += ab_sht_off * teams[j].Players[i].shots_off
			teams[j].Players[i].st_ab += ab_sav * teams[j].Players[i].saves
			teams[j].Players[i].st_ab += ab_concede * teams[j].Players[i].conceded

			// For cards, all abilities are decreased (only St for a GK)
			//
			if teams[j].Players[i].Pos == "GK" {
				teams[j].Players[i].st_ab += ab_yellow * teams[j].Players[i].yellowcards
				teams[j].Players[i].st_ab += ab_red * teams[j].Players[i].redcards
			} else {
				teams[j].Players[i].tk_ab += ab_yellow * teams[j].Players[i].yellowcards
				teams[j].Players[i].ps_ab += ab_yellow * teams[j].Players[i].yellowcards
				teams[j].Players[i].sh_ab += ab_yellow * teams[j].Players[i].yellowcards

				teams[j].Players[i].tk_ab += ab_red * teams[j].Players[i].redcards
				teams[j].Players[i].ps_ab += ab_red * teams[j].Players[i].redcards
				teams[j].Players[i].sh_ab += ab_red * teams[j].Players[i].redcards
			}
		}

		notJ := 1 - j

		// Add random-victory bonuses
		//
		if teams[j].Score > teams[notJ].Score {
			n := 0
			num := 0

			for k := 1; k <= 2; k++ {
				//
				// Find a player to get the increase
				//
				for {
					n = myRandom(numPlayers)
					if teams[j].Players[n].Minutes != 0 && n != num {
						break
					}
				} // while (!teams[j].Players[n].minutes || n == num);

				//
				// Decide the ability which gets the increase
				//
				if teams[j].Players[n].Pos == "GK" {
					teams[j].Players[n].st_ab += ab_victory
				} else {
					teams[j].Players[n].tk_ab += ab_victory
					teams[j].Players[n].ps_ab += ab_victory
					teams[j].Players[n].sh_ab += ab_victory
				}

				num = n
			}
		}

		//
		// Decrease random-defeat bonuses
		//
		if teams[j].Score < teams[notJ].Score {
			n := 0
			num := 0

			for k := 1; k <= 2; k++ {

				//
				// Decide the player to get the decrease
				//
				for {
					n = myRandom(numPlayers)
					if teams[j].Players[n].Minutes != 0 && n != num {
						break
					}

				} // while (!teams[j].Players[n].minutes || n == num);

				//
				// Decide the ability which gets the decrease
				//
				if teams[j].Players[n].Pos == "GK" {
					teams[j].Players[n].st_ab += ab_defeat
				} else {
					teams[j].Players[n].tk_ab += ab_defeat
					teams[j].Players[n].ps_ab += ab_defeat
					teams[j].Players[n].sh_ab += ab_defeat
				}

				num = n
			}
		}

		//
		// Add clean sheet bonus
		//
		if teams[notJ].Score == 0 {
			n := 0

			for {
				n++

				if n >= numPlayers {
					break
				}

				if teams[j].Players[n].Minutes != 0 && teams[j].Players[n].Pos != "GK" {
					break
				}
			} // while (teams[j].Players[n].minutes < 46 || (strcmp(teams[j].Players[n].pos, "GK")));

			if n >= numPlayers {
				n = 1
			}

			teams[j].Players[n].st_ab += ab_cleansheet

			for {
				n = myRandom(numPlayers)

				if teams[j].Players[n].Minutes != 0 && teams[j].Players[n].Pos != "DF" {
					break
				}
			} // while (!teams[j].Players[n].minutes || (strcmp(teams[j].Players[n].pos, "DF")));

			teams[j].Players[n].tk_ab += ab_cleansheet
		}
	}
}

// randomp Generates a random number up to 10000. If the given p is
// less than the generated number, return 1, otherwise return 0
//
// Used to "throw dice" and check if an event with some probability
// happened. p is 0..10000 - for example 2000 means probability 0.2
// So when 2000 is given, this function simulates an event with
// probability 0.2 and tells if it happened (naturally it has
// a prob. of 0.2 to happen)
func randomp(p int) int {
	value := myRandom(10000)

	if value < p {
		return 1
	}
	return 0
}

// myRandom returns a pseudo-random integer between 0 and N-1
func myRandom(n int) int {
	return rnd.Intn(n)
}

// updatePlayersMinuteCount adds one minute to the "minutes played" stats of all currently active
// players in both teams.
func updatePlayersMinuteCount() {
	for j := 0; j <= 1; j++ {
		for i := 0; i < numPlayers; i++ {
			if teams[j].Players[i].Active == 1 {
				teams[j].Players[i].Minutes++
			}
		}
	}
}

// cleanInjCardIndicatorsCalled is called in the beginning of every minute to clean the indicators
// of injuries, yellow and red cards (that are used by conditionals).
func cleanInjCardIndicators() {
	injuredInd[0] = -1
	injuredInd[1] = -1
	yellowCarded[0] = -1
	yellowCarded[1] = -1
	redCarded[0] = -1
	redCarded[1] = -1
}
