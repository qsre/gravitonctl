package cmd

import (
	"gravitonctl/pkg/aws"
	"gravitonctl/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a graviton instance",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.Read()
		if err != nil {
			log.Exit(0)
		}

		if len(args) == 0 {
			log.Error("Please supply a name")
			return
		}
		err = aws.Start(args[0], c.KeyName)
		if err != nil {
			panic(err)
		}

		connectCmd.Run(cmd, args)
	},
}
