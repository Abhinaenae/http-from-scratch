package request

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoodRequestLine(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestGoodRequestLineWithPath(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestGoodPostRequestWithPath(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader("POST /submit HTTP/1.1\r\nHost: localhost:42069\r\nContent-Type: application/json\r\n\r\n{\"key\":\"value\"}"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/submit", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestInvalidNumberOfPartsInRequestLine(t *testing.T) {
	_, err := RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
	require.Error(t, err)
}

func TestInvalidMethodRequestLine(t *testing.T) {
	_, err := RequestFromReader(strings.NewReader("get /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n")) // lowercase method
	require.Error(t, err)
}

func TestInvalidVersionRequestLine(t *testing.T) {
	_, err := RequestFromReader(strings.NewReader("GET /coffee HTTP/2.0\r\nHost: localhost:42069\r\n\r\n")) // Unsupported HTTP version
	require.Error(t, err)
}
