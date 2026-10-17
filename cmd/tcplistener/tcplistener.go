package main

import (
	"fmt"
	"http/internal/request"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("error starting listener")
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("error accepting tcp connection")
		}
		fmt.Println("baka I accepted a connection! ~~~~")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error passing conn to reader")
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
	}
}
