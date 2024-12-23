/*
Package cmd Copyright © 2024 Elladur
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Elladur/tennisaboplaner-go/internal"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tennisaboplaner-go",
	Short: "tennisaboplaner-go helps creating schedules for tennis abos",
	Long: `tennisaboplaner-go automates the process to create an schedule,
	where everybody has prescheduled matches. It can consider various
	requirements from the players. For instance we can exclude certain
	dates for a specific player. This can be used if at the start of 
	the season somebody already knows that he will not be present at a 
	certain week. Additonally we can exclude complete dates from planning
	(e.g. christmas). All these settings can be specified in a 
	settings.json. For an example take look at 
	https://github.com/Elladur/tennisaboplaner-go.`,
	Run: func(cmd *cobra.Command, args []string) {
		content, err := os.ReadFile(cfgFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		var settings internal.SeasonSettings
		err = json.Unmarshal(content, &settings)
		if err != nil {
			fmt.Println(err)
			return
		}
		season, err := internal.CreateSeasonFromSettings(settings)
		if err != nil {
			fmt.Println(err)
			return
		}

		optimizer := internal.Optimizer{Season: &season}
		optimizer.Optimize()
		fmt.Printf("Optimized Schedule and new Score is %.2f\n", internal.GetScore(season.Schedule, season.Players))

		err = season.Export("output")
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "settings.json", "path to config file")
}
