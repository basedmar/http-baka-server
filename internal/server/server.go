package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	state    atomic.Bool
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	serv := Server{state: atomic.Bool{}, listener: l}
	serv.state.Store(true)
	go serv.listen()
	return &serv, nil
}

func (s *Server) Close() error {
	s.state.Store(false)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		con, err := s.listener.Accept()
		if err != nil {
			if x := s.state.Load(); x != true {
				return
			}
			log.Printf("error accepinting con")
			continue
		}
		go s.Handle(con)
	}
}

func (s *Server) Handle(con net.Conn) {
	defer con.Close()
	response := "HTTP/1.1 200 OK\r\n" + // Status line
		"Content-Type: text/plain\r\n" + // Example header
		"\r\n" + // Blank line to separate headers from the body
		"Hello World!\n" // Body
	con.Write([]byte(response))
}
