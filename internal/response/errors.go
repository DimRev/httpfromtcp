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

type ErrorWritingBody struct {
	Err error
}

func (e *ErrorWritingBody) Error() string {
	return fmt.Sprintf("error: writing body: %v", e.Err)
}

type ErrorWritingChunkedBody struct {
	Err error
}

func (e *ErrorWritingChunkedBody) Error() string {
	return fmt.Sprintf("error: writing chunked body: %v", e.Err)
}

type ErrorInvalidWriterState struct {
	CurrentState  writerState
	ExpectedState writerState
}

func (e *ErrorInvalidWriterState) Error() string {
	return fmt.Sprintf("error: invalid writer state: current=%d, expected=%d", e.CurrentState, e.ExpectedState)
}
