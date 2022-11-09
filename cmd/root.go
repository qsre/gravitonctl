package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Execute executes the root command.
func Execute() {
	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(
		startCmd,
		stopCmd,
		terminateCmd,
		deleteCmd,
		connectCmd,
		configureCmd,
		listCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
