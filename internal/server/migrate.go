package server

import (
	"errors"
	"fmt"

	"github.com/DarRo9/Tenders/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func Migrate(cfg *config.PGConfig, log *logrus.Logger) error {
	m, err := migrate.New(
		"file://migrations",
		cfg.Conn)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("no new migrations to apply.")
			return nil 
		}

		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Info("migrations applied successfully.")
	return nil
}
