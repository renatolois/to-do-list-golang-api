package task

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

var ErrNotFound = errors.New("task not found")

type Repository interface {
	Create(t Task) (Task, error)
	GetAll() ([]Task, error)
	GetByID(id string) (Task, error)
	DeleteByID(id string) error
}

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) (*SQLiteRepository, error) {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL,
			completed_at DATETIME,
			created_at DATETIME NOT NULL
		)
	`
	if _, err := db.Exec(query); err != nil {
		return nil, err
	}
	return &SQLiteRepository{db: db}, nil
}

func (r *SQLiteRepository) Create(t Task) (Task, error) {
	query := "INSERT INTO tasks (id, title, completed, completed_at, created_at) VALUES (?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query, t.Id, t.Title, t.Completed, t.CompletedAt, t.CreatedAt)
	if err != nil {
		return Task{}, err
	}
	return t, nil
}

func (r *SQLiteRepository) GetAll() ([]Task, error) {
	query := "SELECT id, title, completed, completed_at, created_at FROM tasks"
	rows, err := r.db.Query(query)
	if err != nil {
		return []Task{}, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.Id, &t.Title, &t.Completed, &t.CompletedAt, &t.CreatedAt); err != nil {
			return []Task{}, err
		}
		tasks = append(tasks, t)
	}
	if tasks == nil {
		return []Task{}, nil
	}
	return tasks, nil
}

func (r *SQLiteRepository) GetByID(id string) (Task, error) {
	query := "SELECT id, title, completed, completed_at, created_at FROM tasks WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var t Task
	err := row.Scan(&t.Id, &t.Title, &t.Completed, &t.CompletedAt, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrNotFound
		}
		return Task{}, err
	}
	return t, nil
}

func (r *SQLiteRepository) DeleteByID(id string) error {
	query := "DELETE FROM tasks WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
