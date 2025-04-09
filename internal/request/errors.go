package request

import "fmt"

type ErrorReadingRequest struct {
	Err error
}

func (e *ErrorReadingRequest) Error() string {
	return fmt.Sprintf("error reading request: %v", e.Err)
}

func (e *ErrorReadingRequest) Unwrap() error {
	return e.Err
}

type ErrorIncompleteRequest struct{}

func (e ErrorIncompleteRequest) Error() string {
	return "incomplete request: missing header termination"
}

type ErrorInvalidRequestLine struct {
	Line string
}

func (e ErrorInvalidRequestLine) Error() string {
	return fmt.Sprintf("invalid request line: %s", e.Line)
}

type ErrorInvalidMethod struct {
	Method string
}

func (e ErrorInvalidMethod) Error() string {
	return fmt.Sprintf("invalid method: %s", e.Method)
}

type ErrorInvalidRequestTarget struct {
	Target string
}

func (e ErrorInvalidRequestTarget) Error() string {
	return fmt.Sprintf("invalid request target: %s", e.Target)
}

type ErrorInvalidHTTPVersion struct {
	Version string
}

func (e ErrorInvalidHTTPVersion) Error() string {
	return fmt.Sprintf("invalid http version: %s", e.Version)
}

type ErrorInvalidHTTPVersionFormat struct {
	Version string
}

func (e ErrorInvalidHTTPVersionFormat) Error() string {
	return fmt.Sprintf("invalid http version format: %s", e.Version)
}

type ErrorInvalidHeaderLine struct {
	HeaderLine string
}

func (e ErrorInvalidHeaderLine) Error() string {
	return fmt.Sprintf("invalid header line: %s", e.HeaderLine)
}
