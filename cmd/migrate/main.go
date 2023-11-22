package main

import (
	"database/sql"
	"github.com/enkodio/pkg-outbox/internal/pkg/config"
	"github.com/enkodio/pkg-outbox/migration/app"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	serviceName = "transaction_outbox_migration"
)

func main() {
	configSettings, err := config.LoadConfigSettingsByPath("configs")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	db, err := sql.Open("pgx", configSettings.PostgresConfigs.GetDSN(serviceName))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	app.Run(db, nil)
}
