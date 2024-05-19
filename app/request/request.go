package request

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Error string
func (e Error) Error() string { return string(e) }

const (
	ErrInvalidRequestLineSize = Error("request line should contain 3 parts")
	ErrFailedToReadLine = Error("failed to read line")
)

type (
	RequestLine struct {
		HTTPMethod  string
		Target      string
		HTTPVersion string
	}

	Request struct {
		RequestLine
		headers     map[string]string
		requestBody string
	}
)

// GET /index.html HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n
// https://datatracker.ietf.org/doc/html/rfc9112#section-2.2
func ParseRequest(data []byte) (Request, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	scanner.Split(scanCRLF)
	
	isEnd := !scanner.Scan()
	if isEnd {
		return Request{}, errors.New("request does not contain a start line")
	}

	requestLine, err := parseRequestLine(scanner.Text())
	if err != nil {
		return Request{}, fmt.Errorf("failed to parse request line: %w", err)
	}

	if requestLine.HTTPMethod == "" {
		requestLine.HTTPMethod = "GET"
	}

	headers := make(map[string]string)

	for scanner.Scan() {
		currentLine := scanner.Text()
		if currentLine == "" || currentLine == "\n" {
			break
		}
		key, val, err :=  parseHeader(scanner, currentLine)
		if err != nil {
			return Request{}, fmt.Errorf("failed to parse header: %w", err)
		}
		headers[key] = val
	}

	var requestBodyBuilder strings.Builder
	for scanner.Scan() {
		requestBodyBuilder.WriteString(scanner.Text())
	}

	return Request{
		RequestLine: requestLine,
		headers:     headers,
		requestBody: requestBodyBuilder.String(),
	}, nil
}

func parseHeader(scanner *bufio.Scanner, line string) (string, string, error) {
	headerLine := strings.SplitN(scanner.Text(), ":", 2)
	if len(headerLine) != 2 {
		return "", "", fmt.Errorf("invalid header: %s", headerLine)
	}
	return headerLine[0], headerLine[1], nil
}

func parseRequestLine(requestLine string) (RequestLine, error) {
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return RequestLine{}, ErrInvalidRequestLineSize
	}

	return RequestLine{HTTPMethod: parts[0], Target: parts[1], HTTPVersion: parts[2]}, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// ScanLines is a split function for a Scanner that returns each line of
// text, stripped of any trailing end-of-line marker. The returned line may
// be empty. The end-of-line marker is one optional carriage return followed
// by one mandatory newline. In regular expression notation, it is `\r?\n`.
// The last non-empty line of input will be returned even if it has no
// newline.
func scanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\r','\n'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}