package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/DimRev/httpfromtcp/internal/request"
	"github.com/DimRev/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	handler  Handler
	listener net.Listener
	closed   atomic.Bool
}

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

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

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

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	w := response.NewWriter()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("Error parsing request:\n- %v\n", err)
		w.WriteStatusLine(response.StatusBadRequest)
		w.WriteHeaders(response.GetDefaultHeaders(len(err.Error())))
		w.WriteBody([]byte(err.Error()))
		w.Write(conn)
		return
	}

	s.handler(w, req)

	err = w.Write(conn)
	if err != nil {
		fmt.Printf("Error writing response:\n- %v\n", err)
		return
	}
}
