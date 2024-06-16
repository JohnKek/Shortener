package sqlite

import (
	"UrlShort/internal/storage"
	"database/sql"
	"fmt"
	"modernc.org/sqlite"
	_ "modernc.org/sqlite" // (1)
)

const (
	NEW    = "storage.sqlite.New"
	SAVE   = "storage.sqlite.SaveURL"
	GET    = "storage.sqlite.GetURL"
	DELETE = "storage.sqlite.DeleteURL"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", NEW, err)
	}
	// TODO Заменить на миграцию
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", NEW, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", NEW, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", SAVE, err)
	}

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		if _, ok := err.(*sqlite.Error); ok {
			return fmt.Errorf("%s: %w", SAVE, storage.ErrURLExists)
		}

		return fmt.Errorf("%s: %w", SAVE, err)
	}
	return nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", GET, err)
	}
	var res string
	err = stmt.QueryRow(alias).Scan(&res)
	if err != nil {
		return "", fmt.Errorf("%s: %w", SAVE, err)
	}
	return res, err

}

func (s *Storage) DeleteURL(alias string) error {
	stmt, err := s.db.Prepare("DELETE url FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", DELETE, err)
	}
	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", SAVE, err)
	}
	return nil
}
