package store

import (
	"database/sql"
	"fmt"
	"jobtrack/models"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func New(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	const schema = `
		CREATE TABLE IF NOT EXISTS applications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			company TEXT NOT NULL,
			role TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TEXT DEFAULT (datetime('now')),
			updated_at TEXT DEFAULT (datetime('now'))
		);`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) ListApplications() ([]models.Application, error) {
	rows, err := s.db.Query(`SELECT id, company, role, status FROM applications`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []models.Application
	for rows.Next() {
		var a models.Application
		if err := rows.Scan(&a.ID, &a.Company, &a.Role, &a.Status); err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, rows.Err()
}

func (s *Store) GetApplication(id int64) (models.Application, error) {
	var a models.Application
	err := s.db.QueryRow(`SELECT id, company, role, status FROM applications WHERE id = ?`, id).Scan(&a.ID, &a.Company, &a.Role, &a.Status)

	if err == sql.ErrNoRows {
		return a, fmt.Errorf("application %d not found", id)
	}
	return a, err
}

func (s *Store) CreateApplication(a *models.Application) error {
	res, err := s.db.Exec(
		`INSERT INTO applications (company, role, status) VALUES (?, ?, ?)`,
		a.Company, a.Role, a.Status,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	a.ID = id
	return nil
}

func (s *Store) DeleteApplication(id int64) error {
	_, err := s.db.Exec(`DELETE FROM applications WHERE id = ?`, id)
	return err
}
