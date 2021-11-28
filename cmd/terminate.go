package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gravitonctl/pkg/aws"
)

var terminateCmd = &cobra.Command{
	Use:   "terminate",
	Short: "Terminates a graviton instance",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error("Please supply a name")
			return
		}
		aws.Terminate(args[0])
		return
	},
}
