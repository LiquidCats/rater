package postgres

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rotisserie/eris"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(conn *pgx.Conn) error {
	sourceDriver, err := iofs.New(migrations, "migrations")
	if err != nil {
		return eris.Wrap(err, "iofs")
	}

	dbConn := stdlib.OpenDB(*conn.Config())
	defer dbConn.Close()

	// Create a new pgx migration driver instance.
	dbDriver, err := pgxmigrate.WithInstance(dbConn, &pgxmigrate.Config{})
	if err != nil {
		return eris.Wrap(err, "pgxmigrate")
	}

	// Create the migrate instance using the source and database drivers.
	m, err := migrate.NewWithInstance(
		"iofs", sourceDriver,
		"pgx", dbDriver,
	)
	if err != nil {
		return eris.Wrap(err, "migrate instance")
	}

	// Run the up migrations.
	if err = m.Up(); err != nil && !eris.Is(err, migrate.ErrNoChange) {
		return eris.Wrap(err, "up")
	}

	return nil
}
