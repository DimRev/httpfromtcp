package server

import "fmt"

type ErrorServerListener struct {
	Err error
}

func (e *ErrorServerListener) Error() string {
	return fmt.Sprintf("error: listening on port: %s", e.Err.Error())
}

type ErrorServerAlreadyClosed struct{}

func (e *ErrorServerAlreadyClosed) Error() string {
	return "error: server already closed"
}

type ErrorServerClose struct {
	Err error
}

func (e *ErrorServerClose) Error() string {
	return fmt.Sprintf("error: closing server: %s", e.Err.Error())
}
