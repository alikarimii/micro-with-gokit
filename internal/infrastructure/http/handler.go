package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alikarimii/micro-with-gokit/internal/endpoints"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	// tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func NewHTTPHandler(endpoints *endpoints.Endpoints, logger log.Logger) *mux.Router {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerBefore(),
	}
	createStudentHandler := httptransport.NewServer(
		endpoints.CreateStudent,
		decodeCreateStudentRequest,
		encodeResponse,
		options...,
	)
	getStudentByIdHandler := httptransport.NewServer(
		endpoints.GetStudentById,
		decodeGetStudentByIdRequest,
		encodeResponse,
		options...,
	)
	getAllStudentsHandler := httptransport.NewServer(
		endpoints.GetAllStudents,
		decodeGetAllStudentsRequest,
		encodeResponse,
		options...,
	)
	deleteStudentHandler := httptransport.NewServer(
		endpoints.DeleteStudent,
		decodeDeleteStudentRequest,
		encodeResponse,
		options...,
	)
	updateStudentHandler := httptransport.NewServer(
		endpoints.UpdateStudent,
		decodeUpdateStudentRequest,
		encodeResponse,
		options...,
	)
	r := mux.NewRouter()
	r.Methods("POST").Path("/v1/student").Handler(createStudentHandler)
	r.Methods("GET").Path("/v1/student/{id}").Handler(getStudentByIdHandler)
	r.Methods("PUT").Path("/v1/student/{id}").Handler(updateStudentHandler)
	r.Methods("DELETE").Path("/v1/student/{id}").Handler(deleteStudentHandler)
	r.Methods("GET").Path("/v1/student").Handler(getAllStudentsHandler)
	r.Handle("/metrics", promhttp.Handler())
	return r
}

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}
func err2code(err error) int {
	// @TODO
	// switch err {
	// case someError:
	// 	return http.StatusBadRequest
	// }
	return http.StatusInternalServerError
}
func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}
