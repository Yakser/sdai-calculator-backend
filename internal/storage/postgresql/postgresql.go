package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"sdai-calculator/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const component = "storage.postgresql.New"

	db, err := sql.Open("pgx", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", component, err)
	}

	return &Storage{
		db: db,
	}, nil

}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const component = "storage.postgresql.SaveURL"

	query, err := s.db.Prepare(`
		insert into urls(url, alias) values($1, $2)
		returning id
	`)
	if err != nil {
		return 0, fmt.Errorf("%v: %w", component, err)
	}
	defer query.Close()

	var id int64
	err = query.QueryRow(urlToSave, alias).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return 0, fmt.Errorf("%v: %w", component, storage.ErrURLAlreadyExists)
			}
		}
		return 0, fmt.Errorf("%v: %w", component, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const component = "storage.postgresql.GetURL"

	query, err := s.db.Prepare(`
		select url from urls where alias = $1
	`)
	if err != nil {
		return "", fmt.Errorf("%v: %w", component, err)
	}

	var url string
	err = query.QueryRow(alias).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%v: %w", component, err)
	}

	return url, nil
}

// todo: add delete method
