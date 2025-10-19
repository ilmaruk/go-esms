package fixtures

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/ilmaruk/go-esms/internal"
	"github.com/ilmaruk/go-esms/internal/database"
)

// const dummy = "DMMY"

var rnd *rand.Rand

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixMicro()))
}

func Create(rootDir, league string, season int) error {
	teams, err := database.LoadClubsByLeague(rootDir, league)
	if err != nil {
		return err
	}

	calendar, err := createFixtures(league, season, teams)
	if err != nil {
		return err
	}

	return database.SaveFixtrues(rootDir, calendar)
}

func createFixtures(league string, season int, teams []internal.Club) (internal.Fixtures, error) {
	num_teams := len(teams)

	if num_teams < 2 {
		return internal.Fixtures{}, fmt.Errorf("two teams or more are needed for a league")
	}

	if num_teams%2 == 1 {
		return internal.Fixtures{}, fmt.Errorf("the number of teams is not even")
		// teams = append(teams, dummy)
		// num_teams++
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
		games[0][i] = teams[i].Code
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

	calendar := internal.Fixtures{
		League: league,
		Season: season,
		Clubs:  map[string]internal.Club{},
		Weeks:  make([]internal.GameWeek, 0, num_weeks_in_round*2),
	}
	for _, t := range teams {
		calendar.Clubs[t.Code] = t
	}
	for week_n := 0; week_n < num_weeks_in_round*2; week_n++ {
		week := internal.GameWeek{
			ID:    week_n + 1,
			Games: make([]internal.Game, 0, num_teams/2),
		}

		for team_n := 0; team_n < num_teams; team_n += 2 {
			var home_team, away_team string
			if week_n < num_weeks_in_round {
				home_team = games[week_n][team_n]
				away_team = games[week_n][team_n+1]
			} else {
				home_team = games[week_n-num_weeks_in_round][team_n+1]
				away_team = games[week_n-num_weeks_in_round][team_n]
			}

			game := internal.Game{
				ID:   makeGameID(league, season, home_team, away_team),
				Home: home_team,
				Away: away_team,
				Seed: rnd.Int63(),
			}

			week.Games = append(week.Games, game)
		}

		calendar.Weeks = append(calendar.Weeks, week)
	}

	return calendar, nil
}

func makeGameID(league string, season int, home, away string) uuid.UUID {
	concat := league + formatSeasonID(season) + home + away
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(concat))
}

func formatSeasonID(id int) string {
	return fmt.Sprintf("%03d", id)
}
