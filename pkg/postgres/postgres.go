package postgres

import (
	"avito-internship/pkg/e"
	"avito-internship/pkg/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PgDatabase struct {
	Pool *pgxpool.Pool
	Dsn  string
}

func NewPgDatabase(pool *pgxpool.Pool, dsn string) *PgDatabase {
	return &PgDatabase{Pool: pool, Dsn: dsn}
}

func Connect() (*PgDatabase, error) {
	const op = "PgDatabase.Connect"
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, e.Wrap(op, err)
	}

	return NewPgDatabase(pool, dsn), nil
}

func (db *PgDatabase) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func (db *PgDatabase) RunMigrations(logger logger.Logger) error {
	const op = "PgDatabase.RunMigrations"

	sqlDb, err := sql.Open("pgx", db.Dsn)
	if err != nil {
		return err
	}
	defer sqlDb.Close()

	driver, err := postgres.WithInstance(sqlDb, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return e.Wrap(op, err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return e.Wrap(op, err)
	}

	logger.Infof("migrations applied successfully")
	return nil
}
