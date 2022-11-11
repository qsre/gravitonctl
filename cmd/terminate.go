package cmd

import (
	"gravitonctl/pkg/aws"
	"gravitonctl/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(terminateCmd, deleteCmd)
}

var terminateCmd = &cobra.Command{
	Use:   "terminate",
	Short: "terminates a graviton instance",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := config.Read()
		if err != nil {
			log.Exit(0)
		}

		if len(args) == 0 {
			log.Error("please supply a name")
			return
		}
		err = aws.Terminate(args[0])
		if err != nil {
			panic(err)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "deletes a graviton instance (terminate alias)",
	Run: func(cmd *cobra.Command, args []string) {
		// delete is an alias of terminate
		terminateCmd.Run(cmd, args)
	},
}
