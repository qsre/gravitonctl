package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "gravitonctl"}

// Execute executes the root command.
func Execute() {
	rootCmd.AddCommand()

	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
