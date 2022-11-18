package application

import (
	"github.com/alikarimii/micro-with-gokit/internal/domain"
)

type (
	CreateStudentRequest struct {
		Student domain.Student
	}
	CreateStudentResponse struct {
		ID  string `json:"id"`
		Err error  `json:"error,omitempty"`
	}
	GetStudentByIdRequest struct {
		ID string `json:"id"`
	}
	GetStudentByIdResponse struct {
		Student interface{} `json:"student,omitempty"`
		Err     string      `json:"error,omitempty"`
	}
	GetAllStudentsRequest struct{}

	GetAllStudentsResponse struct {
		Students interface{} `json:"students,omitempty"`
		Err      string      `json:"error,omitempty"`
	}
	DeleteStudentRequest struct {
		ID string `json:"id"`
	}

	DeleteStudentResponse struct {
		Msg string `json:"msg"`
		Err error  `json:"error,omitempty"`
	}
	UpdateStudentRequest struct {
		Student domain.Student
	}
	UpdateStudentResponse struct {
		Msg string `json:"msg,omitempty"`
		Err error  `json:"error,omitempty"`
	}
)
