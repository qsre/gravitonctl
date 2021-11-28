package main

import (
	"gravitonctl/cmd"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.Info("🌎 gravitonctl starting!")
}

func main() {
	cmd.Execute()
}
