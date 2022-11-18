package application

import (
	"context"

	"github.com/alikarimii/micro-with-gokit/internal/domain"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var _ StudentService = (*studentService)(nil)

func NewBasicService(logger log.Logger, repo StudentRepository) StudentService {
	return &studentService{
		logger,
		repo,
	}
}

type studentService struct {
	logger log.Logger
	repo   StudentRepository
}

// CreateStudent implements StudentService
func (s *studentService) CreateStudent(ctx context.Context, student domain.Student) (string, error) {

	id, _ := uuid.NewUUID()
	customerDetails := domain.Student{
		ID:         id.String(),
		FirstName:  student.FirstName,
		LastName:   student.LastName,
		NationalID: student.NationalID,
	}
	if err := s.repo.CreateStudent(ctx, customerDetails); err != nil {
		level.Error(s.logger).Log("err from repo is ", err)
		return "", err
	}
	return id.String(), nil
}

// DeleteStudent implements StudentService
func (s *studentService) DeleteStudent(ctx context.Context, id string) (string, error) {
	err := s.repo.DeleteStudent(ctx, id)
	if err != nil {
		level.Error(s.logger).Log("err ", err)
		return "", err
	}
	return "ok", nil
}

// GetAllStudents implements StudentService
func (s *studentService) GetAllStudents(ctx context.Context) (interface{}, error) {
	// Use the global TracerProvider.
	tr := otel.Tracer("studentService")
	_, span := tr.Start(ctx, "GetAllStudents")
	span.SetAttributes(attribute.Key("testset").String("this is for test"))
	defer span.End()
	customer, err := s.repo.GetAllStudents(ctx)
	if err != nil {
		level.Error(s.logger).Log("err ", err)
		return nil, err
	}
	return customer, nil
}

// GetStudentById implements StudentService
func (s *studentService) GetStudentById(ctx context.Context, id string) (interface{}, error) {
	customer, err := s.repo.GetStudentById(ctx, id)
	if err != nil {
		level.Error(s.logger).Log("err ", err)
		return nil, err
	}
	return customer, nil
}

// UpdateStudent implements StudentService
func (s *studentService) UpdateStudent(ctx context.Context, student domain.Student) (string, error) {
	msg, err := s.repo.UpdateStudent(ctx, student)
	if err != nil {
		level.Error(s.logger).Log("err from repo is ", err)
		return "", err
	}
	return msg, nil
}
