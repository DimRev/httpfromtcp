package request

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {
	// Group 1: Valid methods
	t.Run("valid methods", func(t *testing.T) {
		// GET method
		r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

		// POST method
		r, err = RequestFromReader(strings.NewReader("POST /submit HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "POST", r.RequestLine.Method)
		assert.Equal(t, "/submit", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

		// PUT method
		r, err = RequestFromReader(strings.NewReader("PUT /update HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "PUT", r.RequestLine.Method)
		assert.Equal(t, "/update", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

		// DELETE method
		r, err = RequestFromReader(strings.NewReader("DELETE /resource HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "DELETE", r.RequestLine.Method)
		assert.Equal(t, "/resource", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	})

	// Group 2: Invalid methods
	t.Run("invalid methods", func(t *testing.T) {
		invalidMethods := []string{
			"PATCH", "HEAD", "OPTIONS", "CONNECT", "TRACE", // Not supported in implementation
			"get", "put", "post", "delete", // Lowercase
			"",                                                  // Empty
			"GETT", "GLAZE", "RIZZ", "PUTT", "POSTT", "DELETET", // Non-existent method
		}

		for _, method := range invalidMethods {
			_, err := RequestFromReader(strings.NewReader(method + " / HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
			require.Error(t, err)
			var errInvalidMethod ErrorInvalidMethod
			require.True(t, errors.As(err, &errInvalidMethod), "Method %s should be invalid", method)
		}

		// Test: Invalid number of parts in request line
		_, err := RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.Error(t, err)
		var errInvalidRequestLine ErrorInvalidRequestLine
		require.True(t, errors.As(err, &errInvalidRequestLine))
	})

	// Group 3: Valid version
	t.Run("valid version", func(t *testing.T) {
		r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
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
			_, err := RequestFromReader(strings.NewReader("GET / " + version + "\r\nHost: localhost:42069\r\n\r\n"))
			require.Error(t, err)
			var errInvalidVersion ErrorInvalidHTTPVersion
			require.True(t, errors.As(err, &errInvalidVersion), "Version %s should be invalid", version)
		}
	})

	// Group 5: Valid targets
	t.Run("valid targets", func(t *testing.T) {
		// Basic path
		r, err := RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)

		// With query parameters
		r, err = RequestFromReader(strings.NewReader("GET /search?q=test&page=1 HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "/search?q=test&page=1", r.RequestLine.RequestTarget)

		// With fragment
		r, err = RequestFromReader(strings.NewReader("GET /page#section HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "/page#section", r.RequestLine.RequestTarget)

		// With URL encoded spaces
		r, err = RequestFromReader(strings.NewReader("GET /path%20with%20spaces HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "/path%20with%20spaces", r.RequestLine.RequestTarget)

		// With multiple path segments
		r, err = RequestFromReader(strings.NewReader("GET /path/with/multiple/segments HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "/path/with/multiple/segments", r.RequestLine.RequestTarget)
	})

	// Group 6: Invalid targets
	t.Run("invalid targets", func(t *testing.T) {
		invalidTargets := []string{
			"",                   // Empty
			"no-leading-slash",   // Missing leading slash
			"http://example.com", // Absolute URL not starting with /
		}

		for _, target := range invalidTargets {
			_, err := RequestFromReader(strings.NewReader("GET " + target + " HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
			require.Error(t, err)
			var errInvalidTarget ErrorInvalidRequestTarget
			require.True(t, errors.As(err, &errInvalidTarget), "Target %s should be invalid", target)
		}
	})

	// Malformed request lines
	t.Run("malformed request lines", func(t *testing.T) {
		malformedLines := []string{
			"GET",                  // Missing target and version
			"GET /",                // Missing version
			"GET / HTTP/1.1 extra", // Too many parts
			" / HTTP/1.1",          // Missing method
			"GET  HTTP/1.1",        // Missing target
		}

		for _, line := range malformedLines {
			_, err := RequestFromReader(strings.NewReader(line + "\r\nHost: localhost:42069\r\n\r\n"))
			require.Error(t, err, "Line %s should cause an error", line)
		}
	})
}
