package application

import (
	"context"

	"github.com/alikarimii/micro-with-gokit/internal/domain"
)

type StudentService interface {
	CreateStudent(ctx context.Context, student domain.Student) (string, error)
	GetStudentById(ctx context.Context, id string) (interface{}, error)
	GetAllStudents(ctx context.Context) (interface{}, error)
	UpdateStudent(ctx context.Context, student domain.Student) (string, error)
	DeleteStudent(ctx context.Context, id string) (string, error)
}

type StudentRepository interface {
	CreateStudent(ctx context.Context, student domain.Student) error
	GetStudentById(ctx context.Context, id string) (interface{}, error)
	GetAllStudents(ctx context.Context) (interface{}, error)
	UpdateStudent(ctx context.Context, student domain.Student) (string, error)
	DeleteStudent(ctx context.Context, id string) error
}
