package config

import (
	log "github.com/sirupsen/logrus"
)

var Config GravitonctlConfig

type GravitonctlConfig struct {
	Region      string
	KeyName     string
	KeyLocation string
}

func init() {
	var err error
	Config, err = Read()
	if err != nil {
		if err.Error() == "config does not exist" {
			log.Error("no config present, run `gravitonctl configure`")
		} else {
			log.Error(err)
		}
	}

}

func NewConfig() GravitonctlConfig {
	var g GravitonctlConfig
	return g
}
