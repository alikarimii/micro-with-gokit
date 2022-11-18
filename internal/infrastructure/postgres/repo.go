package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alikarimii/micro-with-gokit/internal/application"
	"github.com/alikarimii/micro-with-gokit/internal/domain"
	"github.com/go-kit/kit/log"
)

var (
	ErrIdNotFound = errors.New("Id not found")
)

func NewRepo(db *sql.DB, logger log.Logger) application.StudentRepository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "postgres"),
	}
}

var _ application.StudentRepository = (*repo)(nil)

type repo struct {
	db     *sql.DB
	logger log.Logger
}

// CreateStudent implements application.StudentRepository
func (r *repo) CreateStudent(ctx context.Context, student domain.Student) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO student(id, first_name, last_name,national_id) VALUES ($1, $2, $3, $4)", student.ID, student.FirstName, student.LastName, student.NationalID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteStudent implements application.StudentRepository
func (r *repo) DeleteStudent(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM student WHERE id = $1 ", id)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowCnt == 0 {
		return ErrIdNotFound
	}
	return nil
}

// GetAllStudents implements application.StudentRepository
func (r *repo) GetAllStudents(ctx context.Context) (interface{}, error) {
	student := domain.Student{}
	var res []interface{}
	rows, err := r.db.QueryContext(ctx, "SELECT s.id,s.first_name,s.last_name,s.national_id FROM student as s ")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrIdNotFound
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.NationalID)
		res = append([]interface{}{student}, res...)
	}
	return res, nil
}

// GetStudentById implements application.StudentRepository
func (r *repo) GetStudentById(ctx context.Context, id string) (interface{}, error) {
	student := domain.Student{}

	err := r.db.QueryRowContext(ctx, "SELECT s.id,s.first_name,s.last_name,s.national_id FROM student as s where s.id = $1", id).Scan(&student.ID, &student.FirstName, &student.LastName, &student.NationalID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrIdNotFound
		}
		return nil, err
	}
	return student, nil
}

// UpdateStudent implements application.StudentRepository
func (r *repo) UpdateStudent(ctx context.Context, student domain.Student) (string, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE student SET first_name=$1 , last_name=$2 , national_id=$3 WHERE id = $4", student.FirstName, student.LastName, student.NationalID, student.ID)
	if err != nil {
		return "", err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return "", err
	}
	if rowCnt == 0 {
		return "", ErrIdNotFound
	}

	return "updated", err
}
