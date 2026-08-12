package database

import (
	"errors"
	"fmt"
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

	m, err := migrate.New(
		sourceURL,
		databaseURL,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil &&
		!errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
