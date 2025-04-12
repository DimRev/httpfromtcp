package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/DimRev/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, &ErrorServerListener{Err: err}
	}

	s := &Server{
		listener: listener,
	}

	s.closed.Store(false)
	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	if s.closed.Load() {
		return &ErrorServerAlreadyClosed{}
	}

	s.closed.Store(true)
	if s.listener != nil {
		s.listener.Close()
	}

	return nil
}

func (s *Server) listen() {
	for {
		if s.closed.Load() {
			return
		}
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		fmt.Printf("Error writing status line: %v\n", err)
		return
	}

	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil {
		fmt.Printf("Error writing headers: %v\n", err)
		return
	}

}
