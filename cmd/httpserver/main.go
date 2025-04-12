package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DimRev/httpfromtcp/internal/request"
	"github.com/DimRev/httpfromtcp/internal/response"
	"github.com/DimRev/httpfromtcp/internal/server"
)

const PORT = 42069

func main() {
	server, err := server.Serve(PORT, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", PORT)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	handler200(w, req)
}

func handler500(w *response.Writer, req *request.Request) {
	html := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
	h := response.GetDefaultHeaders(len(html))
	h.Replace("Content-Type", "text/html")
	w.WriteStatusLine(response.StatusInternalServerError)
	w.WriteHeaders(h)
	w.WriteBody([]byte(html))
}

func handler400(w *response.Writer, req *request.Request) {
	html := `<html>
	<head>
		<title>400 Bad Request</title>
	</head>
	<body>
		<h1>Bad Request</h1>
		<p>Your request honestly kinda sucked.</p>
	</body>
</html>`
	h := response.GetDefaultHeaders(len(html))
	h.Replace("Content-Type", "text/html")
	w.WriteStatusLine(response.StatusBadRequest)
	w.WriteHeaders(h)
	w.WriteBody([]byte(html))
}

func handler200(w *response.Writer, req *request.Request) {
	html := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
	h := response.GetDefaultHeaders(len(html))
	h.Replace("Content-Type", "text/html")
	w.WriteStatusLine(response.StatusOK)
	w.WriteHeaders(h)
	w.WriteBody([]byte(html))
}
