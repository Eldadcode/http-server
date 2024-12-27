package http

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type HTTPStatusCode string

const (
	HTTP_OK        HTTPStatusCode = "200 OK"
	HTTP_NOT_FOUND HTTPStatusCode = "404 File not found"
)

//go:embed templates/page_not_found.html
var default404Page []byte

type HTTPRequest struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    []byte
}

type HTTPResponse struct {
	Status        HTTPStatusCode
	Server        string
	Date          string
	Content       []byte
	ContentType   string
	ContentLength int
}

var httpMethodHandlers map[string]func(string) (HTTPResponse, error) = map[string]func(string) (HTTPResponse, error){"GET": handleGETRequest}

func parseStartLine(start_line string) (string, string) {
	split_result := strings.Split(start_line, " ")
	http_method, path := split_result[0], split_result[1]
	return http_method, path
}

func parseHTTPRequest(raw_request []byte) (HTTPRequest, error) {
	var request string = string(raw_request)
	var http_request HTTPRequest
	var err error

	split_request := strings.Split(request, "\n")
	method, path := parseStartLine(split_request[0])

	http_request.Method = method
	http_request.Path = path

	return http_request, err
}
func generateHTTPResponse(status HTTPStatusCode, content []byte) HTTPResponse {
	return HTTPResponse{
		Status:        status,
		Server:        "Eldad's GO HTTP Server",
		Date:          time.Now().String(),
		Content:       content,
		ContentType:   "text/html; charset=utf-8",
		ContentLength: len(content),
	}
}
func handleGETRequest(path string) (HTTPResponse, error) {
	var http_response HTTPResponse
	if path == "/" {
		path = "/index.html"
	}

	if _, err := os.Stat(path[1:]); errors.Is(err, os.ErrNotExist) {
		http_response = generateHTTPResponse(HTTP_NOT_FOUND, default404Page)
		return http_response, nil
	}

	file_content, err := os.ReadFile(path[1:])
	if err != nil {
		fmt.Println(err)
		return http_response, err
	}

	http_response = generateHTTPResponse(HTTP_OK, file_content)

	return http_response, err
}

func HandleHTTPRequest(raw_request []byte) (HTTPResponse, error) {
	var http_response HTTPResponse
	http_request, err := parseHTTPRequest(raw_request)
	if err != nil {
		return http_response, err
	}

	handler, ok := httpMethodHandlers[http_request.Method]

	if ok {
		http_response, err = handler(http_request.Path)
	} else {
		log.Printf("Got a %s Request, which is not supported\n", http_request.Method)
		err = fmt.Errorf("unsupported request: %s", http_request.Method)
	}
	return http_response, err

}
