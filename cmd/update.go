package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Elladur/tennisaboplaner-go/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install the latest release",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		info, err := internal.PerformUpdate(ctx)
		if err != nil {
			log.Fatal(err)
			return
		}
		if info.Available {
			fmt.Printf("Update available: %s -> %s\n%s\n", info.Current, info.Latest, info.ReleaseURL)
			return
		}
		if info.Latest == "" {
			fmt.Println("No release found. You are up to date.")
			return
		}
		fmt.Printf("No update needed. You are already on %s (latest: %s).\n", info.Current, info.Latest)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
