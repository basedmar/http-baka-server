package server

import (
	"fmt"
	"http/internal/request"
	"http/internal/response"
	"io"
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

func HandleErr(w io.Writer, hanerr *HandlerError) {
	response.WriteStatusLines(w, response.StatusCode(hanerr.Status))
	response.WriteHeaders(w, response.GetDefaultHeaders(len(hanerr.Msg)))
	w.Write([]byte(hanerr.Msg))
}

func (s *Server) Handle(con net.Conn) {
	defer con.Close()
	req, err := request.RequestFromReader(con)
	if err != nil {
		fmt.Printf("%s", err)
	}
	writer := response.Writer{Statusline: []byte{}}
	s.handle(&writer, req)

	con.Write(writer.Statusline)
	response.WriteHeaders(con, writer.Headers)
	con.Write(writer.Body)

}
