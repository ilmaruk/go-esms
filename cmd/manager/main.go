package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ilmaruk/go-esms/internal"
	"github.com/ilmaruk/go-esms/internal/esms"
	"github.com/ilmaruk/go-esms/internal/roster"
)

var (
	workDir string

	cfg internal.Config

	homeCode string
	awayCode string

	teamCode string
	teamName string
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

func init() {
	// Persistent flags available to all commands
	rootCmd.PersistentFlags().StringVarP(&workDir, "work-dir", "d", ".", "working directory")

	// Local flags for specific commands
	playCmd.Flags().StringVar(&homeCode, "home", "", "home team code")
	playCmd.Flags().StringVar(&awayCode, "away", "", "away team code")
	// addCmd.Flags().StringVarP(&priority, "priority", "p", "medium", "task priority (high, medium, low)")
	// listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "show completed tasks too")

	rosterCreateCmd.Flags().StringVarP(&teamCode, "code", "c", "", "team code")
	rosterCreateCmd.Flags().StringVarP(&teamName, "name", "n", "", "team name")

	// Add subcommands
	rosterCmd.AddCommand(rosterCreateCmd)
	rootCmd.AddCommand(playCmd, rosterCmd)

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
	return esms.Play(workDir, homeCode, awayCode)
}

func createRoster(cmd *cobra.Command, args []string) error {
	return roster.CreateRoster(workDir, teamCode, teamName, cfg.RosterCreator)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
