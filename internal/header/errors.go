package header

type ErrorParsingHeaderMalformed struct {
	Line string
}

func (e *ErrorParsingHeaderMalformed) Error() string {
	return "error: malformed header line: " + e.Line
}

type ErrorParsingHeaderKeyValuePairMissing struct {
	Line string
}

func (e *ErrorParsingHeaderKeyValuePairMissing) Error() string {
	return "error: missing key-value pair in header line: " + e.Line
}

type ErrorParsingHeaderTrailingSpaceInKey struct {
	Line string
}

func (e *ErrorParsingHeaderTrailingSpaceInKey) Error() string {
	return "error: trailing white space in key in header line: " + e.Line
}
