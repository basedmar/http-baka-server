package main

import (
	"crypto/sha256"
	"fmt"
	"http/internal/request"
	"http/internal/response"
	"http/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func errorHandler(w *response.Writer, body string) {
	w.WriteStatusLinesW(3)
	headers := response.GetDefaultHeaders(len(body))
	w.WriteHeadersW(headers)
	w.WriteBody([]byte(body))
}

func videoHandler(w *response.Writer) {
	data, err := os.ReadFile("./assets/vim.mp4")
	if err != nil {
		errorHandler(w, "error getting video")
	}
	w.WriteStatusLinesW(1)
	headers := response.GetDefaultHeaders(len(data))
	headers["content-type"] = "video/mp4"
	w.WriteHeadersW(headers)
	w.WriteBody(data)
}

func Caller400(w *response.Writer) {
	body := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
	w.WriteStatusLinesW(2)
	head := response.GetDefaultHeaders(len(body))
	head["content-type"] = "text/html"
	w.WriteHeadersW(head)
	w.WriteBody([]byte(body))
	return
}

func Caller500(w *response.Writer) {
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
}

func Caller200(w *response.Writer) {
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

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "video") {
		videoHandler(w)
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxy(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		Caller400(w)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		Caller500(w)
		return
	}
	Caller200(w)
}

func proxy(w *response.Writer, req *request.Request) {
	prefix := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	if prefix == "range/4096" {
		headers := response.GetDefaultHeaders(-1)
		headers["Trailer"] = fmt.Sprintf("%s, %s", "X-Content-SHA256", "X-Content-Length")
		headers["Transfer-Encoding"] = "chunked"
		resp, err := http.Get(fmt.Sprintf("https://httpbingo.org/%s", prefix))
		fmt.Println(fmt.Sprintf("https://httpbingo.org/%s", prefix))
		if err != nil {
			fmt.Printf("BAKA I COULDNT READ ITT!!!!!!:() : %s\n", err)
			Caller500(w)
			return
		}
		w.WriteStatusLinesW(1)
		w.WriteHeadersW(headers)
		buf := make([]byte, 1024)
		buf_body := make([]byte, 0)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				fmt.Printf("%d", n)
				w.WriteChunkedBody(buf[:n])
				buf_body = append(buf_body, buf[:n]...)
			}
			if err == io.EOF {
				fmt.Printf("\nFinal body length: %d\n", len(buf_body))
				sha := sha256.Sum256(buf_body)
				headers = make(map[string]string)
				headers["X-Content-Sha256"] = fmt.Sprintf("%x", sha)
				headers["X-Content-Length"] = fmt.Sprintf("%d", len(buf_body))
				w.Write_stream.Write([]byte("0\r\n"))
				w.WriterTrailers(headers)
				break
			}
			if err != nil {
				fmt.Printf("%s", err)
			}
		}
		return
	}
	headers := response.GetDefaultHeaders(-1)
	headers["transfer-encoding"] = "chunked"

	resp, err := http.Get(fmt.Sprintf("https://httpbingo.org/%s", prefix))
	if err != nil {
		fmt.Printf("BAKA I COULDNT READ ITT!!!!!!:() : %s\n", err)
		fmt.Printf(prefix)
		return
	}
	w.WriteStatusLinesW(1)
	w.WriteHeadersW(headers)
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if err == io.EOF {
			w.WriteChunkedBodyDone()
			break
		}
		if err != nil {
			fmt.Printf("%s", err)
		}
		fmt.Printf("%d", n)
		w.WriteChunkedBody(buf[:n])
	}
	return
}

func main() {
	server, err := server.Serve(port, handler)
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
