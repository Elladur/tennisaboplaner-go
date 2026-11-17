/*
Package cmd Copyright © 2024 Elladur
*/
package cmd

import (
	"context"
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
var noUpdateCheck bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "tennisaboplaner-go",
	Version: internal.Version,
	Short:   "tennisaboplaner-go helps creating schedules for tennis abos",
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
		default:
			log.Fatal("invalid log level: ", logLevel)
		}

		// best-effort update check
		if !noUpdateCheck && internal.Version != "dev" {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), internal.DefaultUpdateTimeout)
				defer cancel()
				info, err := internal.CheckForUpdate(ctx)
				if err == nil && info.Available {
					fmt.Fprintf(os.Stderr, "A new release is available: %s -> %s. Run \"tennisaboplaner-go update\" to install it.\nRelease: %s\n", info.Current, info.Latest, info.ReleaseURL)
				}
			}()
		}

		// settings param
		content, err := os.ReadFile(cfgFile)
		if err != nil {
			log.Fatal(err)
		}
		var settings internal.SeasonSettings
		err = json.Unmarshal(content, &settings)
		if err != nil {
			log.Fatal(err)
		}
		// output directory param
		if _, err = os.Stat(outputDirectory); err != nil {
			if !os.IsNotExist(err) {
				log.Fatal(err)
			}
			if err2 := os.MkdirAll(outputDirectory, 0755); err2 != nil {
				log.Fatal(err2)
			}
		}
		log.WithFields(log.Fields{
			"level":  logLevel,
			"config": cfgFile,
			"times":  executionTimes,
			"outDir": outputDirectory,
		}).Info("Running with parameters:")

		err = internal.ExecutePlanerParallel(settings, outputDirectory, executionTimes)
		if err != nil {
			log.Fatal(err)
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
	rootCmd.PersistentFlags().StringVar(&logLevel, "level", "info", "set loglevel (possibilites: debug, info, warn)")
	rootCmd.PersistentFlags().StringVar(&outputDirectory, "outDir", "output", "directory to which the files are exported")
	rootCmd.PersistentFlags().IntVar(&executionTimes, "times", 200, "how many times should be tried to find a optimal schedule")
	rootCmd.PersistentFlags().BoolVar(&noUpdateCheck, "no-update-check", false, "disable automatic update check")
}
