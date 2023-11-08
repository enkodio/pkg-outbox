package app

import (
	"database/sql"
	"flag"
	"github.com/enkodio/pkg-outbox/migration/migrations"
	"github.com/enkodio/pkg-outbox/pkg/environment"
	"github.com/enkodio/pkg-outbox/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	defaultMigrationPath = "."
	defaultCommand       = "status"
)

func setLogger(withLog bool) {
	if withLog {
		goose.SetLogger(logger.GetLogger())
	} else {
		goose.SetLogger(&logrus.Logger{})
	}
}

func getMigrationArgs(migrationArgs map[string]string) (command string, args string, withLog bool) {
	//получаем как переменный окружения при деплое через кубер
	command = environment.GetGooseAction()
	args = environment.GetGooseArgs()
	if command != "" || args != "" {
		return command, args, true
	}

	//получаем как флаги
	if migrationArgs == nil {
		flagCommand := flag.String("c", defaultCommand, "command")
		flagArgs := flag.String("args", "", "args")
		flagWithLog := flag.Bool("l", false, "with log")
		flag.Parse()
		return *flagCommand, *flagArgs, *flagWithLog
	}
	//получаем из переданной мапы(для тестов)
	var ok bool
	if command, ok = migrationArgs["c"]; !ok {
		command = defaultCommand
	}
	args = migrationArgs["args"]

	if _, ok = migrationArgs["l"]; ok {
		withLog = true
	}
	return
}

func saveMigrationsToFile() (string, error) {
	name, value := migrations.GetMigrationCreateOutboxTable()
	file, err := os.Create(name)
	if err != nil {
		return "", errors.Wrap(err, "cant create file")
	}
	_, err = file.Write([]byte(value))
	if err != nil {
		return "", errors.Wrap(err, "cant write migration to file")
	}
	err = file.Close()
	if err != nil {
		logger.GetLogger().WithError(err).Error("cant close file")
	}
	return name, nil
}

func Run(db *sql.DB, serviceName string, migrationArgs map[string]string) {
	name, err := saveMigrationsToFile()
	if err != nil {
		panic(errors.Wrap(err, "cant save migrations to file"))
	}

	command, args, withLog := getMigrationArgs(migrationArgs)
	setLogger(withLog)

	defer func() {
		if err = db.Close(); err != nil {
			panic(errors.Wrap(err, "goose: failed to close DB"))
		}
		if err = os.Remove(name); err != nil {
			logger.GetLogger().WithError(err).Errorf("cant remove migration file %v", name)
		}
	}()
	goose.SetTableName(serviceName)
	if err := goose.Run(command, db, defaultMigrationPath, args); err != nil {
		panic(errors.Wrap(err, "goose: failed run migrations"))
	}
}
