package config

import (
	"fmt"
	"gravitonctl/pkg/aws"
	"os"

	"github.com/mitchellh/go-homedir"
)

func (c GravitonctlConfig) Validate() (err error) {
	// Validate Region
	regions, err := aws.GetRegions()
	if err != nil {
		return err
	}

	var validRegion bool
	for _, region := range regions {
		if c.Region == region {
			validRegion = true
		}
	}

	if !validRegion {
		return fmt.Errorf("region '%s' does not exist", c.Region)
	}

	aws.ReInitWithRegion(c.Region)

	// Validate KeyName
	keys, err := aws.GetKeyNames()
	if err != nil {
		return err
	}

	var validKey bool
	for _, key := range keys {
		if c.KeyName == key {
			validKey = true
		}
	}

	if !validKey {
		return fmt.Errorf("key name '%s' does not exist within your aws account", c.KeyName)
	}

	// Validate KeyPath
	path, err := homedir.Expand(c.KeyLocation)
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	file.Close()

	return nil
}
