package app

import (
	"context"
	"database/sql"
	"flag"
	"github.com/enkodio/pkg-outbox/internal/migration/migrations"
	"github.com/enkodio/pkg-outbox/internal/pkg/environment"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
)

const (
	defaultCommand = "up"
)

func getMigrationArgs(migrationArgs map[string]string) (command string) {
	//получаем как переменный окружения при деплое через кубер
	command = environment.GetGooseAction()
	if command != "" {
		return command
	}

	//получаем как флаги
	if migrationArgs == nil {
		flagCommand := flag.String("c", defaultCommand, "command")
		flag.Parse()
		return *flagCommand
	}
	//получаем из переданной мапы(для тестов)
	var ok bool
	if command, ok = migrationArgs["c"]; !ok {
		command = defaultCommand
	}
	return
}

func Run(db *sql.DB, migrationArgs map[string]string) {
	command := getMigrationArgs(migrationArgs)

	var query string
	switch command {
	case "up":
		query = migrations.GetUpCreateOutboxTable()
	case "down":
		query = migrations.GetDownCreateOutboxTable()
	default:
		panic(errors.New("unknown command"))
	}
	_, err := db.ExecContext(context.Background(), query)
	if err != nil {
		panic(errors.Wrap(err, "cant exec migrations"))
	}
}
