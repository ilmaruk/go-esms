package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ilmaruk/go-esms/internal"
	"github.com/ilmaruk/go-esms/internal/database"
	"github.com/ilmaruk/go-esms/internal/esms"
	"github.com/ilmaruk/go-esms/internal/roster"
	"github.com/ilmaruk/go-esms/internal/teamsheet"
)

var (
	rootDir string

	cfg internal.Config

	homeCode string
	awayCode string

	teamCode string
	teamName string
	avgSkill int

	tactic string
)

var rootCmd = &cobra.Command{
	Use:   "esms",
	Short: "Elecronic Soccer Manager Simulator",
	Long:  `esms is the re-implementation of Eli Benderski's esms in Golang.`,
	// PersistentPreRun: loadTasks,
	// PersistentPostRun: saveTasks,
}

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play a new game",
	// Long:  "Play a new game",
	// Args:  cobra.MinimumNArgs(1),
	RunE: playGame,
}

var rosterCmd = &cobra.Command{
	Use:   "roster",
	Short: "Roster functionalities",
	// Long:  "Play a new game",
	// Args:  cobra.MinimumNArgs(1),
}

var rosterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new roster",
	// Long:  "Play a new game",
	// Args:  cobra.MinimumNArgs(1),
	RunE: createRoster,
}

var teamsheetCmd = &cobra.Command{
	Use:   "teamsheet",
	Short: "Teamsheet functionalities",
	// Long:  "Play a new game",
	// Args:  cobra.MinimumNArgs(1),
}

var teamsheetCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new teamsheet",
	// Long:  "Play a new game",
	// Args:  cobra.MinimumNArgs(1),
	RunE: createTeamsheet,
}

func init() {
	// Persistent flags available to all commands
	rootCmd.PersistentFlags().StringVarP(&rootDir, "root-dir", "d", ".", "root directory")

	// Local flags for specific commands
	playCmd.Flags().StringVar(&homeCode, "home", "", "home team code")
	playCmd.Flags().StringVar(&awayCode, "away", "", "away team code")
	// addCmd.Flags().StringVarP(&priority, "priority", "p", "medium", "task priority (high, medium, low)")
	// listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "show completed tasks too")

	rosterCreateCmd.Flags().StringVarP(&teamCode, "code", "c", "", "team code")
	rosterCreateCmd.Flags().StringVarP(&teamName, "name", "n", "", "team name")
	rosterCreateCmd.Flags().IntVarP(&avgSkill, "skill", "s", 14, "average skill")

	teamsheetCreateCmd.Flags().StringVarP(&teamCode, "code", "c", "", "team code")
	teamsheetCreateCmd.Flags().StringVar(&tactic, "tactic", "442N", "team tactic")

	// Add subcommands
	rosterCmd.AddCommand(rosterCreateCmd)
	teamsheetCmd.AddCommand(teamsheetCreateCmd)
	rootCmd.AddCommand(playCmd, rosterCmd, teamsheetCmd)

	// Setup configuration
	if err := setupConfig(); err != nil {
		panic(err)
	}
}

func setupConfig() error {
	viper.SetConfigName("esms")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	// viper.SetDefault("priority", "medium")
	// viper.SetDefault("file", filepath.Join(os.Getenv("HOME"), ".taskman.json"))

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(&cfg)
}

func playGame(cmd *cobra.Command, args []string) error {
	return esms.Play(rootDir, homeCode, awayCode)
}

func createRoster(cmd *cobra.Command, args []string) error {
	return roster.CreateRoster(rootDir, teamCode, teamName, avgSkill, cfg.RosterCreator)
}

func createTeamsheet(cmd *cobra.Command, args []string) error {
	roster, err := database.LoadRoster(rootDir, teamCode)
	if err != nil {
		return err
	}

	ts, err := teamsheet.CreateTeamsheet(roster, tactic)
	if err != nil {
		return err
	}
	return database.SaveTeamsheet(rootDir, ts)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
