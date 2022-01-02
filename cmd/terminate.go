package cmd

import (
	"gravitonctl/pkg/aws"
	"gravitonctl/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var terminateCmd = &cobra.Command{
	Use:   "terminate",
	Short: "Terminates a graviton instance",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := config.Read()
		if err != nil {
			log.Exit(0)
		}

		if len(args) == 0 {
			log.Error("Please supply a name")
			return
		}
		err = aws.Terminate(args[0])
		if err != nil {
			panic(err)
		}
	},
}
