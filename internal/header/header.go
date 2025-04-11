package header

import (
	"bytes"
	"strings"

	"github.com/DimRev/httpfromtcp/internal/request"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	totalBytesParsed := 0
	for {
		idx := bytes.Index(data, []byte(request.CRLF))
		if idx == -1 {
			return totalBytesParsed, false, nil
		}

		if idx == 0 || strings.TrimSpace(string(data[:idx])) == "" {
			return totalBytesParsed + len(request.CRLF), true, nil
		}

		key, value, err := validateHeader(string(data[:idx]))
		if err != nil {
			return 0, false, err
		}
		h[key] = value
		totalBytesParsed += idx + len(request.CRLF)
		data = data[idx+len(request.CRLF):]
	}
}

func validateHeader(line string) (key, value string, err error) {
	parts := strings.Split(line, ": ")
	if len(parts) != 2 {
		return "", "", &ErrorParsingHeaderKeyValuePairMissing{Line: line}
	}
	if strings.HasSuffix(parts[0], " ") {
		return "", "", &ErrorParsingHeaderTrailingSpaceInKey{Line: line}
	}
	key = strings.TrimSpace(parts[0])
	value = strings.TrimSpace(parts[1])
	return key, value, nil
}
