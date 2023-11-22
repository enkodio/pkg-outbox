package main

import (
	"github.com/enkodio/pkg-outbox/internal/outbox/app"
	"github.com/enkodio/pkg-outbox/internal/pkg/config"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	serviceName = "transaction_outbox"
)

func main() {
	configSettings, err := config.LoadConfigSettingsByPath("internal/cmd/configs")
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}
	app.Run(configSettings, serviceName)
}
