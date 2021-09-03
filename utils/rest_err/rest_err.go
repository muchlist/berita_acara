package rest_err

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//APIError interface untuk mengembalikan error
type APIError interface {
	Message() string
	Status() int
	Error() string
	Causes() []interface{}
}

type apiError struct {
	AStatus  int           `json:"status"`
	AMessage string        `json:"message"`
	AnError  string        `json:"error"`
	ACauses  []interface{} `json:"causes"`
}

func (e *apiError) Status() int {
	return e.AStatus
}

func (e *apiError) Message() string {
	return e.AMessage
}

func (e *apiError) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: [ %v ]",
		e.Message(), e.Status(), e.AnError, e.ACauses)
}

func (e *apiError) Causes() []interface{} {
	return e.ACauses
}

// NewAPIError membuat api error baru dengan mendifinisikan semua isinyas
func NewAPIError(message string, statusCode int, err string, causes []interface{}) APIError {
	return &apiError{
		AStatus:  statusCode,
		AMessage: message,
		AnError:  err,
		ACauses:  causes,
	}
}

// NewAPIErrorFromBytes membuat api error dari []byte
func NewAPIErrorFromBytes(body []byte) (APIError, error) {
	var result apiError
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.New("invalid json body")
	}

	return &result, nil
}

// NewNotFoundError membuat api error ketika objek yang dicari tidak ditemukan
func NewNotFoundError(message string) APIError {
	return &apiError{
		AStatus:  http.StatusNotFound,
		AMessage: message,
		AnError:  "not_found",
		ACauses:  []interface{}{},
	}
}

// NewUnauthorizedError membuat api error user yang tidak diijinkan masuk
func NewUnauthorizedError(message string) APIError {
	return &apiError{
		AStatus:  http.StatusUnauthorized,
		AMessage: message,
		AnError:  "unauthorized",
		ACauses:  []interface{}{},
	}
}

// NewInternalServerError membuat error 500 internal
func NewInternalServerError(message string, err error) APIError {
	result := &apiError{
		AStatus:  http.StatusInternalServerError,
		AMessage: message,
		AnError:  "internal_server_error",
		ACauses:  []interface{}{},
	}
	if err != nil {
		result.ACauses = append(result.ACauses, err.Error())
	}
	return result
}

// NewBadRequestError membuat error jika kesalahan ada pada user
func NewBadRequestError(message string) APIError {
	return &apiError{
		AStatus:  http.StatusBadRequest,
		AMessage: message,
		AnError:  "bad_request",
		ACauses:  []interface{}{},
	}
}