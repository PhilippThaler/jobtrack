package store

import (
	"database/sql"
	"errors"
	"fmt"
	"jobtrack/models"
	"os"

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
		  url TEXT,
			status TEXT NOT NULL,
		  notes TEXT,
		  is_favorite INTEGER DEFAULT 0,
			created_at TEXT DEFAULT (datetime('now')),
			updated_at TEXT DEFAULT (datetime('now'))
		);
		CREATE TABLE IF NOT EXISTS contacts (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  application_id INTEGER REFERENCES applications(id) ON DELETE CASCADE,
		  name TEXT NOT NULL,
		  email TEXT,
		  phone TEXT,
		  role TEXT
		);
		CREATE TABLE IF NOT EXISTS timeline (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  application_id INTEGER REFERENCES applications(id) ON DELETE CASCADE,
		  event TEXT NOT NULL,
		  note TEXT,
		  happened_at TEXT DEFAULT (datetime('now'))
		);
		`

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
	rows, err := s.db.Query(`SELECT id, company, role, status, COALESCE(url, ''), COALESCE(notes, ''), is_favorite, created_at, updated_at FROM applications`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
	}(rows)

	var apps []models.Application
	for rows.Next() {
		var a models.Application
		if err := rows.Scan(&a.ID, &a.Company, &a.Role, &a.Status, &a.URL, &a.Notes, &a.IsFavorite, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, rows.Err()
}

func (s *Store) GetApplication(id int64) (models.Application, error) {
	var a models.Application
	err := s.db.QueryRow(`SELECT id, company, role, status, COALESCE(url, ''), COALESCE(notes, ''), is_favorite, created_at, updated_at FROM applications WHERE id = ?`, id).Scan(&a.ID, &a.Company, &a.Role, &a.Status, &a.URL, &a.Notes, &a.IsFavorite, &a.CreatedAt, &a.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return a, fmt.Errorf("application %d not found", id)
	}
	return a, err
}

func (s *Store) CreateApplication(a *models.Application) error {
	res, err := s.db.Exec(
		`INSERT INTO applications (company, role, status, url, notes, is_favorite) VALUES (?, ?, ?, ?, ?, ?)`,
		a.Company, a.Role, a.Status, a.URL, a.Notes, a.IsFavorite,
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

func (s *Store) UpdateApplication(a *models.Application) error {
	_, err := s.db.Exec(
		`UPDATE applications
		SET company = ?,
		    role = ?,
		    status = ?,
		    url = ?,
		    notes = ?,
		    is_favorite = ?,
		    updated_at = datetime('now')
		WHERE id = ?;`,
		a.Company, a.Role, a.Status, a.URL, a.Notes, a.IsFavorite, a.ID,
	)
	return err
}

func (s *Store) DeleteApplication(id int64) error {
	_, err := s.db.Exec(`DELETE FROM applications WHERE id = ?`, id)
	return err
}

func (s *Store) ToggleFavorite(id int64) error {
	_, err := s.db.Exec(`UPDATE applications SET is_favorite = NOT is_favorite WHERE id = ?`, id)
	return err
}
