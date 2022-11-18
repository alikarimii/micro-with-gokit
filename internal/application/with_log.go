package application

import (
	"context"
	"fmt"
	"time"

	"github.com/alikarimii/micro-with-gokit/internal/domain"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
)

var _ StudentService = (*withLogService)(nil)

func NewWithLogMiddlware(logger log.Logger, repo StudentRepository) StudentService {
	var svc StudentService

	{
		svc = NewBasicService(logger, repo)
		svc = loggingMiddleware(logger)(svc)
	}

	return svc
}

type withLogService struct {
	logger log.Logger
	next   StudentService
}

// CreateStudent implements StudentService
func (s *withLogService) CreateStudent(ctx context.Context, student domain.Student) (string, error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "CreateStudent",
			"input", fmt.Sprintf("%v", student),
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.CreateStudent(ctx, student)
}

// DeleteStudent implements StudentService
func (s *withLogService) DeleteStudent(ctx context.Context, id string) (string, error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "DeleteStudent",
			"input", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.DeleteStudent(ctx, id)
}

// GetAllStudents implements StudentService
func (s *withLogService) GetAllStudents(ctx context.Context) (interface{}, error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "GetAllStudents",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.GetAllStudents(ctx)
}

// GetStudentById implements StudentService
func (s *withLogService) GetStudentById(ctx context.Context, id string) (interface{}, error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "GetStudentById",
			"input", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.GetStudentById(ctx, id)
}

// UpdateStudent implements StudentService
func (s *withLogService) UpdateStudent(ctx context.Context, student domain.Student) (string, error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "UpdateStudent",
			"input", fmt.Sprintf("%v", student),
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.UpdateStudent(ctx, student)
}
