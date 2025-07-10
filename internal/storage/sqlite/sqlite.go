package sqlite

import (
	"database/sql"

	"github.com/atharv3221/rest-api/internal/config"
	_ "modernc.org/sqlite"
)

type SqLite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*SqLite, error) {
	db, err := sql.Open("sqlite", cfg.Storagepath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER)`)

	if err != nil {
		return nil, err
	}

	return &SqLite{
		Db: db,
	}, nil
}

func (s *SqLite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, nil
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}
