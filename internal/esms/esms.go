package esms

import "math/rand"

const (
	numPlayers = 11 // number of players in a team
)

var (
	theConfig = &Config{}

	teams        [2]Team // two teams
	injuredInd   [2]int  // indicators of injured players for both teams
	yellowCarded [2]int  // indicators of yellow carded players for both teams
	redCarded    [2]int  // indicators of red carded players for both teams

	rnd *rand.Rand // random number generator
)

func play() {
}

// func send_off(a, b int) {
//     teams[a].Players[b].yellowcards = 0;
//     teams[a].Players[b].redcards++;
//     teams[a].Players[b].active = 0;

//     if team[a].CurrentGK == b {/* If a GK was sent off */
//         int i = 12, found = 0;

//         if (team[a].substitutions < 3)
//         {
//             while (!found && i <= num_players) /* Look for a keeper on the bench */
//             {
//                 /* If found a keeper */
//                 if (!strcmp(team[a].player[i].pos, "GK") && team[a].player[i].active == 2)
//                 {
//                     int n = 11;

//                     found = 1;

//                     while (team[a].player[n].active != 1) /* Sub him for another player */
//                         n--;
//                     substitute_player(a, n, i, "GK");
//                     team[a].current_gk = i;
//                 }
//                 else
//                 {
//                     found = 0;
//                     i++;
//                 }
//             }

//             if (!found)     /*  If there was no keeper on the bench   */
//             {               /*  Change the position of another player */
//                 int n = 11; /*  (who is on the field) to GK           */

//                 while (team[a].player[n].active != 1)
//                     n--;

//                 change_position(a, n, string("GK"));
//                 team[a].current_gk = n;
//             }
//         }
//         else /* If substitutions >= 3 */
//         {
//             int n = 11;

//             while (team[a].player[n].active != 1)
//                 n--;
//             change_position(a, n, string("GK"));
//             team[a].current_gk = n;
//         }
//     }
// }

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
				} // while (!teams[j].player[n].minutes || n == num);

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

				} // while (!teams[j].player[n].minutes || n == num);

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
			} // while (teams[j].player[n].minutes < 46 || (strcmp(teams[j].player[n].pos, "GK")));

			if n >= numPlayers {
				n = 1
			}

			teams[j].Players[n].st_ab += ab_cleansheet

			for {
				n = myRandom(numPlayers)

				if teams[j].Players[n].Minutes != 0 && teams[j].Players[n].Pos != "DF" {
					break
				}
			} // while (!teams[j].player[n].minutes || (strcmp(teams[j].player[n].pos, "DF")));

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
