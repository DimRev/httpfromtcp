package headers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParse(t *testing.T) {
	// Group 1: Valid headers
	t.Run("valid single header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "localhost:42069", headers["host"])
		assert.Equal(t, len("Host: localhost:42069\r\n"), n)
		assert.False(t, done) // Not done yet, more headers may follow
	})

	t.Run("valid single header with extra whitespace", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069   \r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, "localhost:42069", headers["host"])
		assert.Equal(t, len("Host: localhost:42069   \r\n"), n)
		assert.False(t, done) // Not done yet, more headers may follow
	})

	t.Run("valid 2 headers with existing headers", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069\r\nUser-Agent: curl/7.81.0\r\n\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, "localhost:42069", headers["host"])
		assert.Equal(t, "curl/7.81.0", headers["user-agent"])
		assert.Equal(t, len("Host: localhost:42069\r\nUser-Agent: curl/7.81.0\r\n\r\n"), n)
		assert.True(t, done) // Indicates end of headers
	})

	t.Run("valid done", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("\r\n") // Empty line indicating end of headers
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, 2, n) // Length of CRLF
		assert.True(t, done)  // Indicates end of headers
	})

	// Group 2: Invalid headers
	t.Run("invalid spacing header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("       Host : localhost:42069\r\n")
		n, done, err := headers.Parse(data)
		require.Error(t, err)
		assert.Equal(t, 0, n) // No bytes consumed due to error
		assert.False(t, done) // Not done yet

		var errInvalid *ErrorParsingHeaderTrailingSpaceInKey
		require.True(t, errors.As(err, &errInvalid), "Expected error for invalid spacing header")
	})

	t.Run("missing colon", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host localhost:42069\r\n")
		n, done, err := headers.Parse(data)
		require.Error(t, err)
		assert.Equal(t, 0, n) // No bytes consumed due to error
		assert.False(t, done) // Not done yet

		var errInvalid *ErrorParsingHeaderKeyValuePairMissing
		require.True(t, errors.As(err, &errInvalid), "Expected error for missing colon")
	})

	t.Run("trailing space in key", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host : localhost:42069\r\n")
		n, done, err := headers.Parse(data)
		require.Error(t, err)
		assert.Equal(t, 0, n) // No bytes consumed due to error
		assert.False(t, done) // Not done yet

		var errInvalid *ErrorParsingHeaderTrailingSpaceInKey
		require.True(t, errors.As(err, &errInvalid), "Expected error for trailing space in key")
	})

	t.Run("malformed key", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("H@st: localhost:42069\r\n")
		n, done, err := headers.Parse(data)
		require.Error(t, err)
		assert.Equal(t, 0, n) // No bytes consumed due to error
		assert.False(t, done) // Not done yet

		var errInvalid *ErrorParsingHeaderMalformedKey
		require.True(t, errors.As(err, &errInvalid), "Expected error for malformed key")
	})

	t.Run("empty key", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte(": localhost:42069\r\nHost: localhost:42069\r\n")
		n, done, err := headers.Parse(data)
		require.Error(t, err)
		assert.Equal(t, 0, n) // No bytes consumed due to error
		assert.False(t, done) // Not done yet

		var errInvalid *ErrorParsingHeaderEmptyKey
		require.True(t, errors.As(err, &errInvalid), "Expected error for empty key")
	})

	// Group 3: Edge cases
	t.Run("header with only whitespace", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("   \r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, 2, n) // Length of CRLF
		assert.True(t, done)  // Indicates end of headers
	})

	t.Run("header with trailing whitespace in value", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069   \r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, "localhost:42069", headers["host"])
		assert.Equal(t, len("Host: localhost:42069   \r\n"), n)
		assert.False(t, done) // Not done yet, more headers may follow
	})

	// Group 4: Multiple values for same key
	t.Run("multiple values for same key", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069\r\nHost: localhost:69420\r\na: 1\r\na: 2\r\na: 3\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, "localhost:42069, localhost:69420", headers["host"])
		assert.Equal(t, "1, 2, 3", headers["a"])
		assert.Equal(t, len("Host: localhost:42069\r\nHost: localhost:69420\r\na: 1\r\na: 2\r\na: 3\r\n"), n)
		assert.False(t, done) // Not done yet, more headers may follow
	})
}
