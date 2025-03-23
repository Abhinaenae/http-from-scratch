package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/abhinaenae/http-from-scratch/internal/request"
	"github.com/abhinaenae/http-from-scratch/internal/response"
)

type Server struct {
	handler  Handler
	listener net.Listener
	closed   atomic.Bool
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	msgBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(msgBytes))
	response.WriteHeaders(w, headers)
	w.Write(msgBytes)
}

// Creates a net.Listener and returns a new Server instance. Starts listening for requests inside a goroutine.
func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		handler:  handler,
		listener: listener,
	}

	go s.listen()
	return s, nil
}

// Closes the listener and the server
func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

// Uses a loop to .Accept new connections as they come in, and handles each one in a new goroutine.
// Uses an atomic.Bool to track whether the server is closed or not to ignore connection errors after the server is closed.
func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}

}

//Handles a single connection by writing the following response and then closing the connection:
/*
HTTP/1.1 200 OK
Content-Type: text/plain

Hello World!
*/
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	hErr := s.handler(buf, req)
	if hErr != nil {
		hErr.Write(conn)
		return
	}
	b := buf.Bytes()
	response.WriteStatusLine(conn, response.StatusCodeSuccess)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)
	conn.Write(b)
	return
}
