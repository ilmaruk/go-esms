package roster

import (
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/ilmaruk/go-esms/internal"
	"github.com/ilmaruk/go-esms/internal/database"
)

var (
	rnd *rand.Rand

	// supportedCountries = []string{
	// 	"AR", "AU", "BR", "BG",
	// 	"cam", "HR", "DE", "EN",
	// 	"FR", "DE", "NL", "IE",
	// 	"isr", "IT", "JP", "nig",
	// 	"NO", "ZA", "ES", "US",
	// }
)

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixMicro()))
}

func CreateRoster(workDir, teamCode, teamName string, skill int, cfg internal.RosterCreatorConfig) error {
	var numPlayers = cfg.NumGK + cfg.NumDF + cfg.NumDM + cfg.NumMF + cfg.NumAM + cfg.NumFW

	cfg.AvgMainSkill = skill
	cfg.AvgMidSkill = 11
	cfg.AvgSecondarySkill = skill / 2

	roster := internal.Roster{
		TeamCode: teamCode,
		TeamName: teamName,
		Players:  make(internal.Players, 0, numPlayers),
	}

	// First generate the names for all the players in a single go
	// generator := newParserNameGenerator(parserNameApiKey)
	generator := newRandomUserGenerator()
	persons, err := generator.Generate(numPlayers)
	if err != nil {
		return err
	}

	// Then create the actual players
	for p := 0; p < numPlayers; p++ {
		var player internal.Player

		player.Nationality = persons[p].Country
		player.Name = persons[p].String()
		player.ID = nameToID(player.Name)

		// Age: Varies between 16 and 30
		//
		player.Age = averagedRandom(23, 7)

		// Preferred side: preset probability for each
		//
		temp_rand := unoformRandom(150)

		var temp_side string
		if temp_rand <= 8 {
			temp_side = "RLC"
		} else if temp_rand <= 13 {
			temp_side = "RL"
		} else if temp_rand <= 23 {
			temp_side = "RC"
		} else if temp_rand <= 33 {
			temp_side = "LC"
		} else if temp_rand <= 73 {
			temp_side = "R"
		} else if temp_rand <= 103 {
			temp_side = "L"
		} else {
			temp_side = "C"
		}

		player.PrefSide = temp_side

		half_average_secondary_skill := cfg.AvgSecondarySkill / 2

		// Skills: Depends on the position, first n_goalkeepers
		// will get the highest skill in St, and so on...
		//
		if p < cfg.NumGK {
			player.St = averagedRandomPartDev(cfg.AvgMainSkill, 3)
			player.Tk = averagedRandomPartDev(half_average_secondary_skill, 2)
			player.Ps = averagedRandomPartDev(half_average_secondary_skill, 2)
			player.Sh = averagedRandomPartDev(half_average_secondary_skill, 2)
		} else if p < cfg.NumGK+cfg.NumDF {
			player.Tk = averagedRandomPartDev(cfg.AvgMainSkill, 3)
			player.St = averagedRandomPartDev(half_average_secondary_skill, 2)
			player.Ps = averagedRandomPartDev(cfg.AvgSecondarySkill, 2)
			player.Sh = averagedRandomPartDev(cfg.AvgSecondarySkill, 2)
		} else if p < cfg.NumGK+cfg.NumDF+cfg.NumDM {
			player.Ps = averagedRandomPartDev(cfg.AvgMidSkill, 3)
			player.Tk = averagedRandomPartDev(cfg.AvgMidSkill, 3)
			player.St = averagedRandomPartDev(half_average_secondary_skill, 2)
			player.Sh = averagedRandomPartDev(cfg.AvgSecondarySkill, 2)
		} else if p < cfg.NumGK+cfg.NumDF+cfg.NumDM+cfg.NumMF {
			player.Ps = averagedRandomPartDev(cfg.AvgMainSkill, 3)
			player.St = averagedRandomPartDev(half_average_secondary_skill, 2)
			player.Tk = averagedRandomPartDev(cfg.AvgSecondarySkill, 2)
			player.Sh = averagedRandomPartDev(cfg.AvgSecondarySkill, 2)
		} else if p < cfg.NumGK+cfg.NumDF+cfg.NumDM+cfg.NumMF+cfg.NumAM {
			player.Ps = averagedRandomPartDev(cfg.AvgMidSkill, 3)
			player.Sh = averagedRandomPartDev(cfg.AvgMidSkill, 3)
			player.Tk = averagedRandomPartDev(cfg.AvgSecondarySkill, 2)
			player.St = averagedRandomPartDev(half_average_secondary_skill, 2)
		} else {
			player.Sh = averagedRandomPartDev(cfg.AvgMainSkill, 3)
			player.St = averagedRandomPartDev(half_average_secondary_skill, 2)
			player.Tk = averagedRandomPartDev(cfg.AvgSecondarySkill, 2)
			player.Ps = averagedRandomPartDev(cfg.AvgSecondarySkill, 2)
		}

		// Stamina
		//
		player.Stamina = averagedRandomPartDev(cfg.AvgStamina, 2)

		// Aggression
		//
		player.Ag = float64(averagedRandomPartDev(cfg.AvgAggression, 3))

		// Abilities: set all to 300
		//
		player.ShAb = 300
		player.TkAb = 300
		player.PsAb = 300
		player.ShAb = 300

		// Other stats
		//
		player.Games = 0
		player.Saves = 0
		player.Tackles = 0
		player.KeyPasses = 0
		player.Shots = 0
		player.Goals = 0
		player.Assists = 0
		player.Dp = 0
		player.Injury = 0
		player.Suspension = 0
		player.Fitness = 100

		roster.Players = append(roster.Players, player)
	}

	return database.SaveRoster(workDir, roster)
}

func averagedRandomPartDev(average, div int) int {
	return averagedRandom(float64(average), float64(average)/float64(div))
}

func averagedRandom(average, max_deviation float64) int {
	val := int(rnd.NormFloat64()*max_deviation + average)

	// TODO: this is not great, because values at the edges have a higher probability than expected
	minVal := int(average - max_deviation)
	if val < minVal {
		return minVal
	}
	maxVal := int(average + max_deviation)
	if val > maxVal {
		return maxVal
	}
	return val
}

// Return a pseudo-random integer uniformly distributed
// between 0 and max
func unoformRandom(max int) int {
	return rnd.Intn(max)
	// double d = rand() / (double)RAND_MAX;
	// unsigned u = (unsigned)(d * (max + 1));

	// return (u == max + 1 ? max - 1 : u);
}

func nameToID(name string) string {
	pattern, _ := regexp.Compile("[^a-z]+")
	return pattern.ReplaceAllString(strings.ToLower(name), "_")
}
