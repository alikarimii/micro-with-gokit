package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/alikarimii/micro-with-gokit/internal/application"
	"github.com/gorilla/mux"
)

func decodeCreateStudentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req application.CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Student); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetStudentByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req application.GetStudentByIdRequest
	vars := mux.Vars(r)
	req = application.GetStudentByIdRequest{
		ID: vars["id"],
	}
	return req, nil
}
func decodeGetAllStudentsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req application.GetAllStudentsRequest
	return req, nil
}
func decodeDeleteStudentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req application.DeleteStudentRequest
	vars := mux.Vars(r)
	req = application.DeleteStudentRequest{
		ID: vars["id"],
	}
	return req, nil
}
func decodeUpdateStudentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req application.UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Student); err != nil {
		return nil, err
	}
	vars := mux.Vars(r)
	req.Student.ID = vars["id"]
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
