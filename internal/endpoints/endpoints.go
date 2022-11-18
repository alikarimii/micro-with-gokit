package endpoints

import (
	"context"

	"github.com/alikarimii/micro-with-gokit/internal/application"
	"github.com/alikarimii/micro-with-gokit/pkg"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

func New(svc application.StudentService, logger log.Logger, duration metrics.Histogram) *Endpoints {
	var createStudentEp endpoint.Endpoint
	{
		createStudentEp = makeCreateStudentEndpoint(svc)
		//
		// createStudentEp is limited to 1 request per second with burst of 100 request.
		// Note, rate is defined as a time interval between requests.
		createStudentEp = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(1), 100))(createStudentEp)
		createStudentEp = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(createStudentEp)
		//
		createStudentEp = pkg.LoggingMiddleware(log.With(logger, "method", "CreateStudent"))(createStudentEp)
		createStudentEp = pkg.InstrumentingMiddleware(duration.With("method", "CreateStudent"))(createStudentEp)
	}

	var getStudentByIdEp endpoint.Endpoint
	{
		getStudentByIdEp = makeGetStudentByIdEndpoint(svc)
		getStudentByIdEp = pkg.LoggingMiddleware(log.With(logger, "method", "GetStudentById"))(getStudentByIdEp)
		getStudentByIdEp = pkg.InstrumentingMiddleware(duration.With("method", "GetStudentById"))(getStudentByIdEp)
	}

	var getAllStudentsEp endpoint.Endpoint
	{
		getAllStudentsEp = makeGetAllStudentsEndpoint(svc)
		getAllStudentsEp = pkg.LoggingMiddleware(log.With(logger, "method", "GetAllStudents"))(getAllStudentsEp)
		getAllStudentsEp = pkg.InstrumentingMiddleware(duration.With("method", "GetAllStudents"))(getAllStudentsEp)
	}

	var deleteStudentEp endpoint.Endpoint
	{
		deleteStudentEp = makeDeleteStudentEndpoint(svc)
		//
		// deleteStudentEp is limited to 1 request per second with burst of 1 request.
		deleteStudentEp = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(1), 1))(deleteStudentEp)
		deleteStudentEp = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(deleteStudentEp)
		//
		deleteStudentEp = pkg.LoggingMiddleware(log.With(logger, "method", "DeleteStudent"))(deleteStudentEp)
		deleteStudentEp = pkg.InstrumentingMiddleware(duration.With("method", "DeleteStudent"))(deleteStudentEp)
	}
	var updateStudentEp endpoint.Endpoint
	{
		updateStudentEp = makeUpdateStudentendpoint(svc)
		updateStudentEp = pkg.LoggingMiddleware(log.With(logger, "method", "UpdateStudent"))(updateStudentEp)
		updateStudentEp = pkg.InstrumentingMiddleware(duration.With("method", "UpdateStudent"))(updateStudentEp)
	}

	return &Endpoints{
		CreateStudent:  createStudentEp,
		GetStudentById: getStudentByIdEp,
		GetAllStudents: getAllStudentsEp,
		DeleteStudent:  deleteStudentEp,
		UpdateStudent:  updateStudentEp,
	}
}

type Endpoints struct {
	CreateStudent  endpoint.Endpoint
	GetStudentById endpoint.Endpoint
	GetAllStudents endpoint.Endpoint
	DeleteStudent  endpoint.Endpoint
	UpdateStudent  endpoint.Endpoint
}

func makeCreateStudentEndpoint(s application.StudentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(application.CreateStudentRequest)
		msg, err := s.CreateStudent(ctx, req.Student)
		return application.CreateStudentResponse{ID: msg, Err: err}, nil
	}
}
func makeGetStudentByIdEndpoint(s application.StudentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(application.GetStudentByIdRequest)
		res, err := s.GetStudentById(ctx, req.ID)
		if err != nil {
			return application.GetStudentByIdResponse{Student: res, Err: "Id not found"}, nil
		}
		return application.GetStudentByIdResponse{Student: res, Err: ""}, nil
	}
}
func makeGetAllStudentsEndpoint(s application.StudentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		result, err := s.GetAllStudents(ctx)
		if err != nil {
			return application.GetAllStudentsResponse{Students: result, Err: "no data found"}, nil
		}
		return application.GetAllStudentsResponse{Students: result, Err: ""}, nil
	}
}
func makeDeleteStudentEndpoint(s application.StudentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(application.DeleteStudentRequest)
		msg, err := s.DeleteStudent(ctx, req.ID)
		if err != nil {
			return application.DeleteStudentResponse{Msg: msg, Err: err}, nil
		}
		return application.DeleteStudentResponse{Msg: msg, Err: nil}, nil
	}
}
func makeUpdateStudentendpoint(s application.StudentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(application.UpdateStudentRequest)
		msg, err := s.UpdateStudent(ctx, req.Student)
		return msg, err
	}
}
