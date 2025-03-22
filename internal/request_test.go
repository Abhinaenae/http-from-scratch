package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
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

func TestGoodRequestLine(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestGoodRequestLineWithPath(t *testing.T) {
	reader := &chunkReader{
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestGoodPostRequestWithPath(t *testing.T) {
	reader := &chunkReader{
		data:            "POST /submit HTTP/1.1\r\nHost: localhost:42069\r\nContent-Type: application/json\r\n\r\n{\"key\":\"value\"}",
		numBytesPerRead: 5,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/submit", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestInvalidNumberOfPartsInRequestLine(t *testing.T) {
	reader := &chunkReader{
		data:            "/coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
		numBytesPerRead: 4,
	}
	_, err := RequestFromReader(reader)
	require.Error(t, err)
}

func TestInvalidMethodRequestLine(t *testing.T) {
	reader := &chunkReader{
		data:            "get /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n", // lowercase method
		numBytesPerRead: 2,
	}
	_, err := RequestFromReader(reader)
	require.Error(t, err)
}

func TestInvalidVersionRequestLine(t *testing.T) {
	reader := &chunkReader{
		data:            "GET /coffee HTTP/2.0\r\nHost: localhost:42069\r\n\r\n", // Unsupported HTTP version
		numBytesPerRead: 6,
	}
	_, err := RequestFromReader(reader)
	require.Error(t, err)
}
