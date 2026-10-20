package server

import (
	"fmt"
	"http/internal/request"
	"http/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type HandlerError struct {
	Status int
	Msg    string
}

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	state    atomic.Bool
	listener net.Listener
	handle   Handler
}

func Serve(port int, hander Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	serv := Server{state: atomic.Bool{}, listener: l, handle: hander}
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
	w := response.NewWriter(con)
	req, err := request.RequestFromReader(con)
	if err != nil {
		w.WriteStatusLinesW(2)
		body := []byte(fmt.Sprintf("Error parsing req >_< : %s", err))
		w.WriteHeadersW(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
	}
	s.handle(w, req)

}
