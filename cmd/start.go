package cmd

import (
	"gravitonctl/pkg/aws"
	"gravitonctl/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts a graviton instance",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.Read()
		if err != nil {
			log.Exit(0)
		}

		if len(args) == 0 {
			log.Error("please supply a name")
			return
		}
		err = aws.Start(args[0], c.KeyName)
		if err != nil {
			panic(err)
		}

		connectCmd.Run(cmd, args)
	},
}
