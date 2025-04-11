package request

import (
	"errors"
	"io"
	"testing"

	"github.com/DimRev/httpfromtcp/internal/headers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// chunkReader simulates reading a variable number of bytes from a string.
type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call.
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}

func NewChunkReader(method, target, version string, headers []string, body string, numBytesPerRead int) *chunkReader {
	requestString := generateRequest(method, target, version, body, headers)
	return &chunkReader{
		data:            requestString,
		numBytesPerRead: numBytesPerRead,
	}
}

func generateRequest(method, target, version, body string, headers []string) string {
	// request := method + " " + target + " " + version + "\r\n"
	request := ""
	if method != "" {
		request += method + " "
	}
	if target != "" {
		request += target + " "
	}
	if version != "" {
		request += version
	}
	request += "\r\n"

	for _, header := range headers {
		request += header + "\r\n"
	}

	request += "\r\n" + body
	return request
}

func TestRequestLineParse(t *testing.T) {
	// Group 1: Valid methods
	t.Run("valid methods", func(t *testing.T) {
		headers := []string{
			"Host: localhost:42069",
			"User-Agent: curl/7.81.0",
			"Accept: */*",
		}
		reader := NewChunkReader("GET", "/", "HTTP/1.1", headers, "", 3)

		r, err := RequestFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
		assert.Equal(t, "localhost:42069", r.Headers["host"])
		assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
		assert.Equal(t, "*/*", r.Headers["accept"])
		assert.Equal(t, []byte{}, r.Body)

		// POST method
		reader = NewChunkReader("POST", "/submit", "HTTP/1.1", headers, "", 3)
		r, err = RequestFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "POST", r.RequestLine.Method)
		assert.Equal(t, "/submit", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
		assert.Equal(t, "localhost:42069", r.Headers["host"])
		assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
		assert.Equal(t, "*/*", r.Headers["accept"])
		assert.Equal(t, []byte{}, r.Body)

		// PUT method
		reader = NewChunkReader("PUT", "/update", "HTTP/1.1", headers, "", 3)
		r, err = RequestFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "PUT", r.RequestLine.Method)
		assert.Equal(t, "/update", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
		assert.Equal(t, "localhost:42069", r.Headers["host"])
		assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
		assert.Equal(t, "*/*", r.Headers["accept"])
		assert.Equal(t, []byte{}, r.Body)

		// DELETE method
		reader = NewChunkReader("DELETE", "/resource", "HTTP/1.1", headers, "", 3)
		r, err = RequestFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "DELETE", r.RequestLine.Method)
		assert.Equal(t, "/resource", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
		assert.Equal(t, "localhost:42069", r.Headers["host"])
		assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
		assert.Equal(t, "*/*", r.Headers["accept"])
		assert.Equal(t, []byte{}, r.Body)
	})

	// Group 2: Invalid methods
	t.Run("invalid methods", func(t *testing.T) {
		invalidMethods := []string{
			"PATCH", "HEAD", "OPTIONS", "CONNECT", "TRACE", // Not supported in implementation
			"get", "put", "post", "delete", // Lowercase
			"GETT", "GLAZE", "RIZZ", "PUTT", "POSTT", "DELETET", // Non-existent method
		}

		for _, method := range invalidMethods {
			_, err := RequestFromReader(NewChunkReader(method, "/", "HTTP/1.1", []string{"Host: localhost:42069"}, "", 3))
			require.Error(t, err)
			var errInvalidMethod *ErrorParsingRequestInvalidMethod
			require.True(t, errors.As(err, &errInvalidMethod), "Method %s should be invalid", method)
		}

		// Test: Invalid number of parts in request line
		_, err := RequestFromReader(NewChunkReader("", "/coffee", "HTTP/1.1", []string{"Host: localhost:42069"}, "", 3))
		require.Error(t, err)
		var errInvalidRequestLine *ErrorParsingRequestLineMalformed
		require.True(t, errors.As(err, &errInvalidRequestLine))
	})

	// Group 3: Valid version
	t.Run("valid version", func(t *testing.T) {
		r, err := RequestFromReader(NewChunkReader("GET", "/", "HTTP/1.1", []string{"Host: localhost:42069"}, "", 3))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	})

	// Group 4: Invalid versions
	t.Run("invalid versions", func(t *testing.T) {
		invalidVersions := []string{
			"HTTP/1.0", // Not 1.1
			"HTTP/2.0", // Not 1.1
			"HTTP/0.9", // Not 1.1
			"HTTP/",    // No version number
			"HTTP",     // Missing slash and version
			"http/1.1", // Lowercase
			"1.1",      // Missing HTTP/
			"",         // Empty
		}

		for _, version := range invalidVersions {
			_, err := RequestFromReader(NewChunkReader("GET", "/", version, []string{"Host: localhost:42069"}, "", 3))
			require.Error(t, err)
			var errInvalidVersion *ErrorParsingRequestInvalidVersion
			require.True(t, errors.As(err, &errInvalidVersion), "Version %s should be invalid", version)
		}
	})

	// Group 5: Valid targets
	t.Run("valid targets", func(t *testing.T) {
		// Basic path
		r, err := RequestFromReader(NewChunkReader("GET", "/coffee", "HTTP/1.1", []string{"Host: localhost:42069"}, "", 3))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)

		// With query parameters
		r, err = RequestFromReader(NewChunkReader("GET", "/search?q=test&page=1", "HTTP/1.1", []string{"Host: localhost:42069"}, "", 3))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "/search?q=test&page=1", r.RequestLine.RequestTarget)

		// With fragment
		r, err = RequestFromReader(NewChunkReader("GET", "/page#section", "HTTP/1.1", []string{"Host: localhost:42069"}, "", 3))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "/page#section", r.RequestLine.RequestTarget)

		// With URL encoded spaces
		r, err = RequestFromReader(NewChunkReader("GET", "/path%20with%20spaces", "HTTP/1.1", []string{"Host: localhost:42069"}, "", 3))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "/path%20with%20spaces", r.RequestLine.RequestTarget)

		// With multiple path segments
		r, err = RequestFromReader(NewChunkReader("GET", "/path/with/multiple/segments", "HTTP/1.1", []string{"Host: localhost:42069"}, "", 3))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "/path/with/multiple/segments", r.RequestLine.RequestTarget)
	})

	// Group 6: Invalid targets
	t.Run("invalid targets", func(t *testing.T) {
		invalidTargets := []string{
			"no-leading-slash",   // Missing leading slash
			"http://example.com", // Absolute URL not starting with /
		}

		for _, target := range invalidTargets {
			_, err := RequestFromReader(NewChunkReader("GET", target, "HTTP/1.1", []string{"Host: localhost:42069"}, "", 3))
			require.Error(t, err)
			var errInvalidTarget *ErrorParsingRequestInvalidTarget
			require.True(t, errors.As(err, &errInvalidTarget), "Target %s should be invalid", target)
		}
	})

	// Group 7: Malformed request lines
	t.Run("malformed request lines", func(t *testing.T) {
		malformedLines := []string{
			"GET",                  // Missing target and version
			"GET /",                // Missing version
			"GET / HTTP/1.1 extra", // Too many parts
			" / HTTP/1.1",          // Missing method
			"GET  HTTP/1.1",        // Missing target
		}

		for _, line := range malformedLines {
			_, err := RequestFromReader(NewChunkReader("GET", line, "HTTP/1.1", []string{"Host: localhost:42069"}, "", 3))
			require.Error(t, err, "Line %s should cause an error", line)
		}
	})

	// Group 8: Variant chunk sizes
	t.Run("variant chunk sizes", func(t *testing.T) {
		chunkSizes := []int{
			1, 10, 100, 1000, 10000,
		}
		for _, chunkSize := range chunkSizes {
			r, err := RequestFromReader(NewChunkReader(
				"GET",
				"/coffee",
				"HTTP/1.1",
				[]string{
					"Host: localhost:42069",
					"Transfer-Encoding: chunked",
					"User-Agent: curl/7.81.0",
					"Accept: */*",
					"",
					"Connection: close",
				},
				"",
				chunkSize,
			),
			)
			require.NoError(t, err)
			require.NotNil(t, r)
		}
	})

	// Group 9: Malformed headers
	t.Run("malformed headers", func(t *testing.T) {
		_, err := RequestFromReader(NewChunkReader("GET", "/", "HTTP/1.1", []string{"Host : localhost:42069"}, "", 3))
		require.Error(t, err)
		var errInvalidHeader *headers.ErrorParsingHeaderTrailingSpaceInKey
		require.True(t, errors.As(err, &errInvalidHeader), "Expected error for malformed header")

		_, err = RequestFromReader(NewChunkReader("GET", "/", "HTTP/1.1", []string{"Host localhost:42069"}, "", 3))
		require.Error(t, err)
		var errInvalidHeaderKV *headers.ErrorParsingHeaderKeyValuePairMissing
		require.True(t, errors.As(err, &errInvalidHeaderKV), "Expected error for malformed header")

		_, err = RequestFromReader(NewChunkReader("GET", "/", "HTTP/1.1", []string{"H@st: localhost:42069"}, "", 3))
		require.Error(t, err)
		var errInvalidHeaderMalformedKey *headers.ErrorParsingHeaderMalformedKey
		require.True(t, errors.As(err, &errInvalidHeaderMalformedKey), "Expected error for malformed header")
	})

	// Group 10: Request with body
	// Group 10: Request with body
	t.Run("Request with body", func(t *testing.T) {
		body := "Hello World!\n"
		r, err := RequestFromReader(NewChunkReader("POST", "/test", "HTTP/1.1", []string{"Content-Length: 13"}, body, 3))
		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, "13", r.Headers["content-length"])
		assert.Equal(t, []byte(body), r.Body)

		_, err = RequestFromReader(NewChunkReader("POST", "/test", "HTTP/1.1", []string{"Content-Length: 3"}, "Hello World!\n", 10))
		require.Error(t, err)

		_, err = RequestFromReader(NewChunkReader("POST", "/test", "HTTP/1.1", []string{"Content-Length: 15"}, "Hello World!\n", 10))
		require.Error(t, err)
	})
}
