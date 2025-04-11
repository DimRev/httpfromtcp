package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

var validKeyChars = map[rune]bool{
	'a': true, 'b': true, 'c': true, 'd': true, 'e': true,
	'f': true, 'g': true, 'h': true, 'i': true, 'j': true,
	'k': true, 'l': true, 'm': true, 'n': true, 'o': true,
	'p': true, 'q': true, 'r': true, 's': true, 't': true,
	'u': true, 'v': true, 'w': true, 'x': true, 'y': true,
	'z': true,

	'A': true, 'B': true, 'C': true, 'D': true, 'E': true,
	'F': true, 'G': true, 'H': true, 'I': true, 'J': true,
	'K': true, 'L': true, 'M': true, 'N': true, 'O': true,
	'P': true, 'Q': true, 'R': true, 'S': true, 'T': true,
	'U': true, 'V': true, 'W': true, 'X': true, 'Y': true,
	'Z': true,

	'0': true, '1': true, '2': true, '3': true, '4': true,
	'5': true, '6': true, '7': true, '8': true, '9': true,

	'!': true, '#': true, '$': true, '%': true, '&': true,
	'\'': true, '*': true, '+': true, '-': true, '.': true,
	'^': true, '_': true, '`': true, '|': true, '~': true,
}

const CRLF = "\r\n"

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	totalBytesParsed := 0
	for {
		idx := bytes.Index(data, []byte(CRLF))
		if idx == -1 {
			return totalBytesParsed, false, nil
		}

		if idx == 0 || strings.TrimSpace(string(data[:idx])) == "" {
			return totalBytesParsed + len(CRLF), true, nil
		}

		key, value, err := validateHeader(string(data[:idx]))
		if err != nil {
			return 0, false, err
		}

		curr, exists := h[strings.ToLower(key)]
		if exists {
			h[strings.ToLower(key)] = fmt.Sprintf("%s, %s", curr, value)
		} else {
			h[strings.ToLower(key)] = value
		}

		totalBytesParsed += idx + len(CRLF)
		data = data[idx+len(CRLF):]
	}
}

func (h Headers) Get(key string) string {
	val, exists := h[strings.ToLower(key)]
	if !exists {
		return ""
	}
	return val
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
	if len(key) == 0 {
		return "", "", &ErrorParsingHeaderEmptyKey{Line: line}
	}
	if !isValidHeaderKey(key) {
		return "", "", &ErrorParsingHeaderMalformedKey{Line: line}
	}

	value = strings.TrimSpace(parts[1])
	return key, value, nil
}

func isValidHeaderKey(key string) bool {
	for _, c := range key {
		if !validKeyChars[c] {
			return false
		}
	}
	return true
}
