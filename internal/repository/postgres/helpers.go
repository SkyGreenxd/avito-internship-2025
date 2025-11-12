package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func postgresDuplicate(err, errIsExists error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return errIsExists
	}
	return err
}

func checkGetQueryResult(err, errNotFound error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return errNotFound
	}

	return err
}
