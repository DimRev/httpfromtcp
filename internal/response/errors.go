package response

import "fmt"

type ErrorInvalidStatusCode struct {
	StatusCode int
}

func (e *ErrorInvalidStatusCode) Error() string {
	return fmt.Sprintf("error: invalid status code: %d", e.StatusCode)
}

type ErrorWritingStatusLine struct {
	Err error
}

func (e *ErrorWritingStatusLine) Error() string {
	return fmt.Sprintf("error: writing status line: %v", e.Err)
}

type ErrorWritingHeaders struct {
	Err error
}

func (e *ErrorWritingHeaders) Error() string {
	return fmt.Sprintf("error: writing headers: %v", e.Err)
}
