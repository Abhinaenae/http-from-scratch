package request

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	//parse request data
	req, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading http request %v", err)
	}
	lines := strings.Split(string(req), "\r\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty request")
	}

	//extract requestline data
	requestLineData := lines[0]

	parts := strings.Split(requestLineData, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line format")
	}

	method, requestTarget, httpVersion := parts[0], parts[1], parts[2]

	//make sure method is capital
	validMethod := regexp.MustCompile(`^[A-Z]+$`)
	if !validMethod.MatchString(method) {
		return nil, fmt.Errorf("invalid HTTP method: %s", method)
	}

	//check http version
	if httpVersion != "HTTP/1.1" {
		return nil, fmt.Errorf("unsupported HTTP version: %s", httpVersion)
	}
	httpVersion = "1.1"

	//Set request line
	requestLine := RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   httpVersion,
	}

	return &Request{
		RequestLine: requestLine,
	}, nil

}
