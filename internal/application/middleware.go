package application

import (
	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type middleware func(StudentService) StudentService

// LoggingMiddleware takes a logger as a dependency
// and returns a service Middleware.
func loggingMiddleware(logger log.Logger) middleware {
	return func(next StudentService) StudentService {
		return &withLogService{logger, next}
	}
}
