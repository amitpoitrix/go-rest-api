package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/amitpoitrix/students-api/internal/config"
	"github.com/amitpoitrix/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

/*
We need to install Go Sqlite Driver and add it in import but don't have to explicitly use in the code
use "_" to ignore the unused error as we're using it indirectly
*/

/* Implementing Storage interface */
type Sqlite struct {
	Db *sql.DB
}

/*
As we don't have contructor in Go so by convention we create New() that act as contructor or initialise
initial values
*/
func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	// Above 3 ? is placeholder so as to prevent SQL Injection
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ?")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	/* Now we've to serialize the data being fetch from DB */
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

func (s *Sqlite) GetAllStudents() ([]types.Student, error) {
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

func (s *Sqlite) ModifyStudentById(id int64, modifyStudent types.Student) (types.Student, error) {
	// 1. first fetch current student whether it exists or not
	currentStudent, err := s.GetStudentById(id)
	if err != nil {
		return types.Student{}, err
	}

	// 2. now dynamically forming update query
	setClause := []string{}
	args := []interface{}{}

	if modifyStudent.Name != "" {
		setClause = append(setClause, "name = ?")
		args = append(args, modifyStudent.Name)
	}

	if modifyStudent.Email != "" {
		setClause = append(setClause, "email = ?")
		args = append(args, modifyStudent.Email)
	}

	if modifyStudent.Age != 0 {
		setClause = append(setClause, "age = ?")
		args = append(args, modifyStudent.Age)
	}

	if len(setClause) == 0 {
		slog.Info("no new data to update for", slog.String("studentId", fmt.Sprint(id)))
		return currentStudent, nil
	}

	// 3. now forming the UPDATE query to update new data
	query := fmt.Sprintf("UPDATE students SET %s WHERE id = ?", strings.Join(setClause, ", "))
	args = append(args, id)

	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		return types.Student{}, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return types.Student{}, err
	}

	if rowsAffected == 0 {
		return types.Student{}, fmt.Errorf("no student found with id: %d", id)
	}

	return s.GetStudentById(id)
}

func (s *Sqlite) DeleteStudentById(id int64) (types.Student, error) {
	// 1. first fetch current student whether it exists or not
	currentStudent, err := s.GetStudentById(id)
	if err != nil {
		return types.Student{}, err
	}

	// 2. now forming the DELETE query to update new data
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return types.Student{}, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return types.Student{}, err
	}

	if rowsAffected == 0 {
		return types.Student{}, fmt.Errorf("no student found with id: %d", id)
	}

	return currentStudent, nil
}
