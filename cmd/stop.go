package cmd

import (
	"gravitonctl/pkg/aws"
	"gravitonctl/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stops a graviton instance",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.Read()
		if err != nil {
			log.Exit(0)
		}

		if len(args) == 0 {
			log.Error("please supply a name")
			return
		}
		err = aws.Stop(args[0], c.KeyName)
		if err != nil {
			panic(err)
		}
	},
}
