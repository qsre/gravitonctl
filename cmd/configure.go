package cmd

import (
	"fmt"
	"gravitonctl/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "configures graviton ctl",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("configuring gravitonctl")

		// REGION
		fmt.Printf("AWS Region [eu-central-1]: ")
		var region string
		l, err := fmt.Scanln(&region)
		if err != nil {
			if err.Error() != "unexpected newline" {
				log.Error(err)
			}
		}

		if l == 0 {
			region = "eu-central-1"
		}

		// KEY NAME
		fmt.Printf("default key name [ie: myKey]: ")
		var keyName string
		l, err = fmt.Scanln(&keyName)
		if err != nil {
			if err.Error() != "unexpected newline" {
				log.Error(err)
			}
		}

		if l == 0 {
			log.Info("key name cannot be blank, exiting ...")
			return
		}

		// KEY LOCATION
		fmt.Printf("default key location [ie: ~/myKey.pem]: ")
		var keyLocation string
		l, err = fmt.Scanln(&keyLocation)
		if err != nil {
			if err.Error() != "unexpected newline" {
				log.Error(err)
			}
		}

		if l == 0 {
			log.Error("key location cannot be blank, exiting ...")
			return
		}

		config := config.GravitonctlConfig{
			Region:      region,
			KeyName:     keyName,
			KeyLocation: keyLocation,
		}

		err = config.Validate()
		if err != nil {
			log.Error(err)
			return
		}

		err = config.Write()
		if err != nil {
			log.Error(err)
			return
		}

		log.Println("gravitonctl configured sucessfully!")

	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "configures graviton ctl (configure alias)",
	Run: func(cmd *cobra.Command, args []string) {
		configureCmd.Run(cmd, args)
	},
}
