/*
Package cmd Copyright © 2024 Elladur
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Elladur/tennisaboplaner-go/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfgFile string
var logLevel string
var outputDirectory string
var executionTimes int

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
		// loglevel switch
		switch logLevel {
		case "info":
			log.SetLevel(log.InfoLevel)
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "debug":
			log.SetLevel(log.DebugLevel)
		}

		// settings param
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
		// output directory param
		_, err = os.Stat(outputDirectory)
		if err != nil {
			if err2 := os.MkdirAll(outputDirectory, 0777); err2 != nil {
				log.Fatal(err)
			}
		}
		log.WithFields(log.Fields{
			"level":  logLevel,
			"config": cfgFile,
			"times":  executionTimes,
			"outDir": outputDirectory,
		}).Info("Running with parameters:")

		internal.ExecutePlanerParallel(settings, outputDirectory, executionTimes)
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
	rootCmd.PersistentFlags().StringVar(&logLevel, "level", "info", "set loglevel (possibilites: debug, info, warn)")
	rootCmd.PersistentFlags().StringVar(&outputDirectory, "outDir", "output", "directory to which the files are exported")
	rootCmd.PersistentFlags().IntVar(&executionTimes, "times", 100, "how many times should be tried to find a optimal schedule")
}
