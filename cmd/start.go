package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gravitonctl/pkg/aws"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a graviton instance",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error("Please supply a name")
			return
		}
		aws.Start(args[0])

		connectCmd.Run(cmd, args)
		return
	},
}
