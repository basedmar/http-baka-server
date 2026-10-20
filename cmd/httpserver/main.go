package main

import (
	"fmt"
	"http/internal/request"
	"http/internal/response"
	"http/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func yourproblemHandler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		body := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
		err := w.WriteStatusLinesW(2)
		if err != nil {
			fmt.Printf("%s", err)
		}
		head := response.GetDefaultHeaders(len(body))
		head["content-type"] = "text/html"
		err = w.WriteHeadersW(head)
		if err != nil {
			fmt.Printf("%s", err)
		}
		err = w.WriteBody([]byte(body))
		if err != nil {
			fmt.Printf("%s", err)
		}
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		body := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
		w.WriteStatusLinesW(3)
		head := response.GetDefaultHeaders(len(body))
		head["content-type"] = "text/html"
		w.WriteHeadersW(head)
		w.WriteBody([]byte(body))
		return
	} else {
		body := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
		w.WriteStatusLinesW(1)
		head := response.GetDefaultHeaders(len(body))
		head["content-type"] = "text/html"
		w.WriteHeadersW(head)
		w.WriteBody([]byte(body))
		return
	}
}

func main() {
	server, err := server.Serve(port, yourproblemHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
