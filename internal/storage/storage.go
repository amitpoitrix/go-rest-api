package storage

import "github.com/amitpoitrix/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetAllStudents() ([]types.Student, error)
	ModifyStudentById(id int64, updateStudent types.Student) (types.Student, error)
	DeleteStudentById(id int64) (types.Student, error)
}
