package sqlstore

import (
	"context"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type sqlstore struct {
	*sqlx.DB

	Driver string
}

func New(driver, dsn string) (*sqlstore, error) {
	dbx, err := sqlx.Connect(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlstore: failed to connect to database")
	}

	db := &sqlstore{dbx, driver}

	// run pending migrations
	if err = db.AutoMigrate(); err != nil {
		return db, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return db, err
	}

	return db, nil
}

func (s *sqlstore) AutoMigrate() error {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationsFS,
		Root:       "migrations",
	}

	migrate.SetSchema("migrations")

	n, err := migrate.Exec(s.DB.DB, s.Driver, migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("sqlstore: failed to apply migrations - %v", err)
	}

	if n > 0 {
		log.Printf("applied %d database migrations", n)
	}

	return nil
}

func (s *sqlstore) Close() error {
	return s.DB.Close()
}
