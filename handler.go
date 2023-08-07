package twomesmanualserver

import (
	"log"
	"net/http"
)

// A HandlerError contains information about an error that occured inside a Handler.
type HandlerError struct {
	Err  error
	Code int
}

func (e HandlerError) Error() string {
	return e.Err.Error()
}

// Create a new HandlerError with error and code.
func NewHandlerError(err error, code int) *HandlerError {
	return &HandlerError{
		Err:  err,
		Code: code,
	}
}

// A Handler is an http.HandlerFunc that can return an error.
type Handler func(w http.ResponseWriter, r *http.Request) error

// Implement the http.Handler interface.
func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)

	if err != nil {
		log.Println(err)

		handlerErr, ok := err.(*HandlerError)
		if !ok {
			HTTPError(w, http.StatusInternalServerError)
			return
		}

		HTTPError(w, handlerErr.Code)
	}
}

// Send an error to the HTTP client with the status text and code.
func HTTPError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
