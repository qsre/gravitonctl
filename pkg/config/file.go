package config

import (
	"encoding/json"
	"fmt"
	"gravitonctl/pkg/aws"
	"os"

	"github.com/mitchellh/go-homedir"
)

func (c GravitonctlConfig) Write() (err error) {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	// remove relative path from keyLocation before writing
	keyLocation, err := homedir.Expand(c.KeyLocation)
	if err != nil {
		return err
	}

	c.KeyLocation = keyLocation

	dir := fmt.Sprintf("%s/.gravitonctl", home)
	fullPath := fmt.Sprintf("%s/config", dir)

	b, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		return err
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(fullPath, b, 0755)
	if err != nil {
		return err
	}

	return err
}

func Read() (c GravitonctlConfig, err error) {
	home, err := homedir.Dir()
	if err != nil {
		return c, err
	}

	dir := fmt.Sprintf("%s/.gravitonctl/config", home)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return c, fmt.Errorf("config does not exist")
	}

	data, err := os.ReadFile(dir)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}

	aws.ReInitWithRegion(c.Region)

	return c, err
}
