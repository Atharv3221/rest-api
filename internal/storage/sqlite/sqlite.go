package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/atharv3221/rest-api/internal/config"
	"github.com/atharv3221/rest-api/internal/types"
	_ "modernc.org/sqlite"
)

type SqLite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*SqLite, error) {
	slog.Info("Starting storage initialization")
	db, err := sql.Open("sqlite", cfg.Storagepath)
	if err != nil {
		slog.Error("while opening database")
		return nil, err
	}

	slog.Info("Storage initialized sucessfully", slog.String("env", cfg.Env))

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER)`)

	if err != nil {
		slog.Error("query error")
		return nil, err
	}

	slog.Info("Successfully created table students")

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

func (s *SqLite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()
	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil
}

func (s *SqLite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err

		}
		students = append(students, student)
	}
	return students, nil
}

func (s *SqLite) DeleteById(id int64) error {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		slog.Error("Internal Error Querey preparation")
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(id)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("student not found with id", slog.Int64("id", id))
		}
		return err
	}
	return nil
}

func (s *SqLite) UpdateStudent(name string, email string, age int, id int64) error {
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		slog.Error("Internal server error query")
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(name, email, age, id)
	if err != nil {
		slog.Error("Internal error", slog.String("Error", err.Error()))
	}
	return nil
}
