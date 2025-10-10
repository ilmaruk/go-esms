package fixtures

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const dummy = "DMMY"

var rnd *rand.Rand

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixMicro()))
}

func CreateFixtures(rootDir string, teams []string) error {
	num_teams := len(teams)

	if num_teams < 2 {
		return fmt.Errorf("two teams or more are needed for a league")
	}

	if num_teams%2 == 1 {
		teams = append(teams, dummy)
		num_teams++
	}

	// Initialize the games vector
	//
	num_weeks_in_round := num_teams - 1
	// vector<string> empty_round(num_teams);
	// vector<vector<string>> games(num_weeks_in_round, empty_round);
	games := make([][]string, num_weeks_in_round)
	for i := range games {
		games[i] = make([]string, num_teams)
	}

	// Initialize 1st week (1st vs. 2nd, 3rd vs. 4th, etc...)
	//
	for i := 0; i < num_teams; i++ {
		games[0][i] = teams[i]
	}

	// Create a round of games
	//
	if num_teams > 2 {
		for week_n := 0; week_n < num_weeks_in_round-1; week_n++ {
			// Each week is built from the previous week, using the
			// algorithm described at the top of this file.
			//
			for team_n := 1; team_n < num_teams-1; team_n++ {
				if team_n%2 == 1 {
					games[week_n+1][team_n+2] = games[week_n][team_n]
				} else {
					games[week_n+1][team_n-2] = games[week_n][team_n]
				}
			}

			// Special rotation around the first team (which doesn't move)
			//
			games[week_n+1][0] = games[week_n][0]
			games[week_n+1][1] = games[week_n][2]
			games[week_n+1][num_teams-2] = games[week_n][num_teams-1]
		}
	}

	// Calibrate home/away so that every team playes home-away-home-away...
	// (very approximately: better for large leagues, worse for small ones).
	//
	// This is done by swapping all teams' home/away every other week.
	//
	for week_n := 1; week_n < num_weeks_in_round; week_n += 2 {
		for team_n := 0; team_n < num_teams; team_n += 2 {
			// Swap
			games[week_n][team_n], games[week_n][team_n+1] = games[week_n][team_n+1], games[week_n][team_n]
		}
	}

	fh, err := os.Create(filepath.Join(rootDir, "data", "fixtures.txt"))
	if err != nil {
		return err
	}
	defer fh.Close()

	// Fixtures calendar = {};
	for week_n := 0; week_n < num_weeks_in_round*2; week_n++ {
		// FixturesWeek week;

		for team_n := 0; team_n < num_teams; team_n += 2 {
			var home_team, away_team string
			if week_n < num_weeks_in_round {
				home_team = games[week_n][team_n]
				away_team = games[week_n][team_n+1]
			} else {
				home_team = games[week_n-num_weeks_in_round][team_n+1]
				away_team = games[week_n-num_weeks_in_round][team_n]
			}

			fmt.Fprintln(fh, home_team, away_team)

			// week.push_back({home_team, away_team});
		}

		// calendar.weeks.push_back(week);
	}

	return nil
}
