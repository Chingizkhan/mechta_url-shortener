package app

import (
	"errors"
	"github.com/Chingizkhan/url-shortener/config"
	"github.com/Chingizkhan/url-shortener/db"
	"github.com/Chingizkhan/url-shortener/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"os"
)

func Migrate(cfg *config.Config, l *logger.Logger) {
	driver, err := iofs.New(db.Migrations, "migrations")
	if err != nil {
		l.Error("Migrate: iofs.New error:", logger.Err(err))
		os.Exit(1)
	}
	m, err := migrate.NewWithSourceInstance(
		"iofs",
		driver,
		cfg.DSN())
	if err != nil {
		l.Error("Migrate: postgres is trying to connect:", logger.Err(err))
		os.Exit(1)
	}

	err = m.Up()
	defer m.Close()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			l.Info("Migrate: no change")
			return
		}

		l.Error("Migrate: up error:", logger.Err(err))
		os.Exit(1)
	}

	l.Info("Migrate: up success")
}
