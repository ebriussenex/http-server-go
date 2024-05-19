package status

import (
	"fmt"

	"github.com/codecrafters-io/http-server-starter-go/app/format"
)

const (
	codeNotFound = 404
	textNotFound = "Not Found"

	codeOk = 200
	textOk = "OK"
)



type StatusCode struct {
	text string
	code uint16
}

var StatusOK = StatusCode{
	text: textOk, code: codeOk,
}

var StatusNotFound = StatusCode{
	text: textNotFound, code: codeNotFound,
}

func StatusLine(protocolVersion string, statusCode StatusCode) string {
	return fmt.Sprintf("%s %d %s", protocolVersion, statusCode.code, statusCode.text) + format.Crlf + format.Crlf
}
