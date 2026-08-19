package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL string) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}

	migrations := filepath.Join(root, "migrations")
	sourceURL := fmt.Sprintf("file://%s", migrations)

	log.Printf("migration database: %s", databaseURL)
	log.Printf("migration source: %s", sourceURL)

	m, err := migrate.New(
		sourceURL,
		databaseURL,
	)
	if err != nil {
		return err
	}

	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("read migration version: %w", err)
	}

	log.Printf(
		"migration state: version=%d dirty=%t",
		version,
		dirty,
	)

	if err := m.Up(); err != nil &&
		!errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
