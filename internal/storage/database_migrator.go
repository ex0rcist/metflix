package storage

import (
	"errors"
	"time"

	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/utils"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DatabaseMigrator struct {
	dsn     string
	retries int
	source  string

	migrator *migrate.Migrate
	err      error
}

func NewDatabaseMigrator(dsn string, source string, retries int) DatabaseMigrator {
	return DatabaseMigrator{dsn: dsn, source: source, retries: retries}
}

func (m DatabaseMigrator) Run() error {
	m.err = utils.NewRetrier(
		func() error {
			logging.LogInfo("migrations: connecting to " + m.dsn)

			migrator, err := migrate.New(m.source, m.dsn)
			m.migrator = migrator

			if err != nil {
				logging.LogError(err)
			}

			return err
		},
		func(err error) bool {
			return true
		},
		[]time.Duration{
			1 * time.Second,
			3 * time.Second,
			5 * time.Second,
		},
	).Run()

	if m.err != nil {
		return m.err
	}

	m.err = m.migrator.Up()

	defer func() {
		srcErr, dbErr := m.migrator.Close()

		if srcErr != nil {
			logging.LogError(srcErr, "failed closing migrator", srcErr.Error())
		}

		if dbErr != nil {
			logging.LogError(dbErr, "failed closing migrator", dbErr.Error())
		}
	}()

	if m.err == nil {
		logging.LogInfo("migrations: success")
		return nil
	}

	if errors.Is(m.err, migrate.ErrNoChange) {
		logging.LogInfo("migrations: no change")
		return nil
	}

	return m.err
}
