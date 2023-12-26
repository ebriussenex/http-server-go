package main

import (
	"fmt"
	"net"
	"os"
)

const (
	crlf = "\r\n"
)

type StatusCode struct {
	text string
	code uint8
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("accepted connection")

		var result []byte

		if _, err := connection.Read(result); err != nil {
			fmt.Println("Failed to read from connection")
		}

		statusLine := formStatusLine("HTTP v1.1", StatusCode{"OK", 200})
		connection.Write([]byte(statusLine))
	}
}

func formStatusLine(protocolVersion string, statusCode StatusCode) string {
	return fmt.Sprintf("%s %d %s", protocolVersion, statusCode.code, statusCode.text) + crlf + crlf
}
