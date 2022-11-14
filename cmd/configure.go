package cmd

import (
	"fmt"
	"gravitonctl/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configureCmd, configCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "configures graviton ctl",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("configuring gravitonctl")

		c, err := config.Read()
		if err != nil {
			log.Error(err)
		}

		// REGION
		if c.Region != "" {
			fmt.Printf("AWS Region [%s]: ", c.Region)
		} else {
			fmt.Printf("AWS Region [ie: eu-central-1]: ")
		}

		l, err := fmt.Scanln(&c.Region)
		if err != nil {
			if err.Error() != "unexpected newline" {
				log.Error(err)
			}
		}

		if l == 0 && c.Region == "" {
			log.Error("region cannot be blank, exiting ...")
			return
		}

		// KEY NAME
		if c.KeyName != "" {
			fmt.Printf("default key name [%s]:", c.KeyName)
		} else {
			fmt.Printf("default key name [ie: myKey]: ")
		}

		l, err = fmt.Scanln(&c.KeyName)
		if err != nil {
			if err.Error() != "unexpected newline" {
				log.Error(err)
			}
		}

		if l == 0 && c.KeyName == "" {
			log.Error("key name cannot be blank, exiting ...")
			return
		}

		// KEY LOCATION
		if c.KeyLocation != "" {
			fmt.Printf("default key location [%s]: ", c.KeyLocation)
		} else {
			fmt.Printf("default key location [ie: ~/myKey.pem]: ")
		}

		l, err = fmt.Scanln(&c.KeyLocation)
		if err != nil {
			if err.Error() != "unexpected newline" {
				log.Error(err)
			}
		}

		if l == 0 && c.KeyLocation == "" {
			log.Error("key location cannot be blank, exiting ...")
			return
		}

		err = c.Validate()
		if err != nil {
			log.Error(err)
			return
		}

		err = c.Write()
		if err != nil {
			log.Error(err)
			return
		}

		log.Info("gravitonctl configured sucessfully!")

	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "configures graviton ctl (configure alias)",
	Run: func(cmd *cobra.Command, args []string) {
		configureCmd.Run(cmd, args)
	},
}
