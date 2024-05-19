package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/request"
	"github.com/codecrafters-io/http-server-starter-go/app/status"
)

type HTTPHandlerFunc func(net.Conn, []byte) error

type Server struct {
	listener net.Listener
	handlers map[string]map[string]HTTPHandlerFunc
}

func (s *Server) registerHandler(httpMethod string, target string, handler HTTPHandlerFunc) {
	if s.handlers[httpMethod] == nil {
		s.handlers[httpMethod] = map[string]HTTPHandlerFunc{}
	}

	s.handlers[httpMethod][target] = handler
}

func (s *Server) handleNotFound(connection net.Conn) error {
	if _, err := connection.Write([]byte(status.StatusLine("HTTP/1.1", status.StatusNotFound))); err != nil {
		return err
	}
	connection.Close()
	return nil
}

func (s *Server) handle(connection net.Conn, data []byte) error {
	fmt.Println("handling request")

	request, err := request.ParseRequest(data)
	if err != nil {
		return fmt.Errorf("failed to parse request: %w", err)
	}

	if handler, ok := s.handlers[request.HTTPMethod][request.Target]; ok {
		return handler(connection, data)
	}

	if request.HTTPMethod == "GET" || request.HTTPMethod == "" {
		return s.handleNotFound(connection)
	}

	if err := connection.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return fmt.Errorf("no handler for %s %s", request.HTTPMethod, request.Target)
}

func (s *Server) Serve() error {
	for {
		var (
			err    error
			amount int
			connection net.Conn
		)

		data := make([]byte, 1024)
		if connection, err = s.listener.Accept(); err != nil {
			return fmt.Errorf("error accepting connection: %w", err)
		}

		fmt.Println("accepted connection!")

		if amount, err = connection.Read(data); err != nil {
			return fmt.Errorf("failure while reading from connection: %w", err)
		}

		fmt.Printf("read from connection: %s, %d\n", string(data), amount)

		if amount > 0 {
			fmt.Printf("served incoming data: %s, %d\n", string(data), amount)

			if err := s.handle(connection, data); err != nil {
				return err
			}
		}
	}
}

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	fmt.Println("server started")

	server := Server{
		l,
		map[string]map[string]HTTPHandlerFunc{},
	}

	server.registerHandler("GET", "/", func(connection net.Conn, _ []byte) error {
		fmt.Println("writing started")
		connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		connection.Close()
		return nil
	})

	if err := server.Serve(); err != nil {
		log.Println(err.Error())
	}
	defer l.Close()
}
