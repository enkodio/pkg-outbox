package main

import (
	"database/sql"
	"github.com/enkodio/pkg-outbox/internal/migration/app"
	"github.com/enkodio/pkg-outbox/internal/pkg/config"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	serviceName = "transaction_outbox_migration"
)

func main() {
	configSettings, err := config.LoadConfigSettingsByPath("internal/cmd/configs")
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
