package progress

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}

	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS progress (
			item_id TEXT PRIMARY KEY,
			completed INTEGER NOT NULL DEFAULT 0,
			updated_at TEXT NOT NULL DEFAULT (datetime('now'))
		)
	`)
	return err
}

func (s *Store) IsComplete(itemID string) (bool, error) {
	var completed int
	err := s.db.QueryRow(`SELECT completed FROM progress WHERE item_id = ?`, itemID).Scan(&completed)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return completed == 1, nil
}

func (s *Store) MarkComplete(itemID string) error {
	_, err := s.db.Exec(`
		INSERT INTO progress (item_id, completed, updated_at)
		VALUES (?, 1, datetime('now'))
		ON CONFLICT(item_id) DO UPDATE SET completed = 1, updated_at = datetime('now')
	`, itemID)
	return err
}

func (s *Store) All() (map[string]bool, error) {
	rows, err := s.db.Query(`SELECT item_id, completed FROM progress WHERE completed = 1`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make(map[string]bool)
	for rows.Next() {
		var id string
		var completed int
		if err := rows.Scan(&id, &completed); err != nil {
			return nil, err
		}
		if completed == 1 {
			items[id] = true
		}
	}
	return items, rows.Err()
}
