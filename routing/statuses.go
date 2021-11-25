package routing

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type ResponseStatus struct {
	Err        error  `json:"-"`
	StatusCode int    `json:"-"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}

var (
	StatusOK            = &ResponseStatus{StatusCode: 200, Message: "OK"}
	ErrMethodNotAllowed = &ResponseStatus{StatusCode: 405, Message: "Method not allowed"}
	ErrNotFound         = &ResponseStatus{StatusCode: 404, Message: "Resource not found"}
	ErrBadRequest       = &ResponseStatus{StatusCode: 400, Message: "Bad request"}
)

func (e *ResponseStatus) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func ErrorRenderer(err error) *ResponseStatus {
	return &ResponseStatus{
		Err:        err,
		StatusCode: 400,
		StatusText: "Bad request",
		Message:    err.Error(),
	}
}

func ServerErrorRenderer(err error) *ResponseStatus {
	return &ResponseStatus{
		Err:        err,
		StatusCode: 500,
		StatusText: "Internal server error",
		Message:    err.Error(),
	}
}

func JSONMarshal(data interface{}) ([]byte) {
	e, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	return e
}

func HTTPResponse(w http.ResponseWriter, statusCode int, resp []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)
}
