package main

import (
	"gravitonctl/cmd"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("🌎 gravitonctl")
	cmd.Execute()
}
