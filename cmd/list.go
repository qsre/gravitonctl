package cmd

import (
	"gravitonctl/pkg/aws"
	"gravitonctl/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all graviton instances launched by gravitonctl",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := config.Read()
		if err != nil {
			log.Exit(0)
		}

		instances, err := aws.DescribeAllRunningInstances()
		if err != nil {
			log.Error(err)
		}

		if len(instances) == 0 {
			log.Println("no running instances")
		}

		for _, instance := range instances {
			var name string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
				}
			}

			log.Printf("%s - %s", *instance.PublicIpAddress, name)
		}
	},
}
