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

	cfg_n_gk = 3
	cfg_n_df = 8
	cfg_n_dm = 3
	cfg_n_mf = 8
	cfg_n_am = 3
	cfg_n_fw = 5

	cfg_average_stamina         = 0
	cfg_average_aggression      = 0
	cfg_average_main_skill      = 4
	cfg_average_mid_skill       = 1
	cfg_average_secondary_skill = 7
)

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixMicro()))
}

func CreateRoster(workDir, teamCode, teamName string) error {
	var numPlayers = 25

	roster := internal.Roster{
		TeamCode: teamCode,
		TeamName: teamName,
		Players:  make(internal.Players, 0, numPlayers),
	}

	// First generate the names for all the players in a single go
	persons, err := generatePersons(numPlayers)
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
		player.Age = averaged_random(23, 7)

		// Preferred side: preset probability for each
		//
		temp_rand := uniform_random(150)

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

		half_average_secondary_skill := cfg_average_secondary_skill / 2

		// Skills: Depends on the position, first n_goalkeepers
		// will get the highest skill in St, and so on...
		//
		if p <= cfg_n_gk {
			player.St = averaged_random_part_dev(cfg_average_main_skill, 3)
			player.Tk = averaged_random_part_dev(half_average_secondary_skill, 2)
			player.Ps = averaged_random_part_dev(half_average_secondary_skill, 2)
			player.Sh = averaged_random_part_dev(half_average_secondary_skill, 2)
		} else if p <= cfg_n_gk+cfg_n_df {
			player.Tk = averaged_random_part_dev(cfg_average_main_skill, 3)
			player.St = averaged_random_part_dev(half_average_secondary_skill, 2)
			player.Ps = averaged_random_part_dev(cfg_average_secondary_skill, 2)
			player.Sh = averaged_random_part_dev(cfg_average_secondary_skill, 2)
		} else if p <= cfg_n_gk+cfg_n_df+cfg_n_dm {
			player.Ps = averaged_random_part_dev(cfg_average_mid_skill, 3)
			player.Tk = averaged_random_part_dev(cfg_average_mid_skill, 3)
			player.St = averaged_random_part_dev(half_average_secondary_skill, 2)
			player.Sh = averaged_random_part_dev(cfg_average_secondary_skill, 2)
		} else if p <= cfg_n_gk+cfg_n_df+cfg_n_dm+cfg_n_mf {
			player.Ps = averaged_random_part_dev(cfg_average_main_skill, 3)
			player.St = averaged_random_part_dev(half_average_secondary_skill, 2)
			player.Tk = averaged_random_part_dev(cfg_average_secondary_skill, 2)
			player.Sh = averaged_random_part_dev(cfg_average_secondary_skill, 2)
		} else if p <= cfg_n_gk+cfg_n_df+cfg_n_dm+cfg_n_mf+cfg_n_am {
			player.Ps = averaged_random_part_dev(cfg_average_mid_skill, 3)
			player.Sh = averaged_random_part_dev(cfg_average_mid_skill, 3)
			player.Tk = averaged_random_part_dev(cfg_average_secondary_skill, 2)
			player.St = averaged_random_part_dev(half_average_secondary_skill, 2)
		} else {
			player.Sh = averaged_random_part_dev(cfg_average_main_skill, 3)
			player.St = averaged_random_part_dev(half_average_secondary_skill, 2)
			player.Tk = averaged_random_part_dev(cfg_average_secondary_skill, 2)
			player.Ps = averaged_random_part_dev(cfg_average_secondary_skill, 2)
		}

		// Stamina
		//
		player.Stamina = averaged_random_part_dev(cfg_average_stamina, 2)

		// Aggression
		//
		player.Ag = float64(averaged_random_part_dev(cfg_average_aggression, 3))

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

func more_st(p1, p2 internal.Player) bool {
	return p1.St > p2.St
}

func more_tk(p1, p2 internal.Player) bool {
	return p1.Tk > p2.Tk
}

func more_ps(p1, p2 internal.Player) bool {
	return p1.Ps > p2.Ps
}

func more_sh(p1, p2 internal.Player) bool {
	return p1.Sh > p2.Sh
}

func averaged_random_part_dev(average, div int) int {
	return averaged_random(float64(average), float64(average)/float64(div))
}

func averaged_random(average, max_deviation float64) int {
	return int(rnd.NormFloat64()*max_deviation + average)
}

// Given a string with comma separated values (like "a,cd,k")
// returns a random value.
func rand_elem(csv string) string {
	elems := strings.Split(csv, ",")
	return elems[uniform_random(len(elems))]
}

// Return a pseudo-random integer uniformly distributed
// between 0 and max
func uniform_random(max int) int {
	return rnd.Intn(max)
	// double d = rand() / (double)RAND_MAX;
	// unsigned u = (unsigned)(d * (max + 1));

	// return (u == max + 1 ? max - 1 : u);
}

// Throws a bet with probability prob of success. Returns
// true if succeeded.
func throw_with_prob(prob int) bool {
	a_throw := 1 + uniform_random(99)

	return prob >= a_throw
}

func nameToID(name string) string {
	pattern, _ := regexp.Compile("[^a-z]+")
	return pattern.ReplaceAllString(strings.ToLower(name), "_")
}
