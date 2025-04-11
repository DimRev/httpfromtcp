package request

import "fmt"

type ErrorIncompleteRequest struct{}

func (e *ErrorIncompleteRequest) Error() string {
	return "error: incomplete request"
}

type ErrorUnexpectedReadError struct {
	Err error
}

func (e *ErrorUnexpectedReadError) Error() string {
	return fmt.Sprintf("error: reading request: %s", e.Err.Error())
}

type ErrorParsingUnknownState struct {
	State requestState
}

func (e *ErrorParsingUnknownState) Error() string {
	return fmt.Sprintf("error: unknown state %d", e.State)
}

type ErrorParsingTryingToReadAfterDone struct{}

func (e *ErrorParsingTryingToReadAfterDone) Error() string {
	return "error: trying to read after done"
}

type ErrorParsingRequestLineMalformed struct {
	Line string
}

func (e *ErrorParsingRequestLineMalformed) Error() string {
	return fmt.Sprintf("error: malformed request line: %s", e.Line)
}

type ErrorParsingRequestInvalidMethod struct {
	Method string
}

func (e *ErrorParsingRequestInvalidMethod) Error() string {
	return fmt.Sprintf("error: invalid method: %s", e.Method)
}

type ErrorParsingRequestInvalidTarget struct {
	Target string
}

func (e *ErrorParsingRequestInvalidTarget) Error() string {
	return fmt.Sprintf("error: invalid target: %s", e.Target)
}

type ErrorParsingRequestInvalidVersion struct {
	Version string
}

func (e *ErrorParsingRequestInvalidVersion) Error() string {
	return fmt.Sprintf("error: invalid version: %s", e.Version)
}
