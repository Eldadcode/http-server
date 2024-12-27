package http

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
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
	Version string
}

func (r HTTPRequest) String() string {
	return fmt.Sprintf("%s %s %s", r.Method, r.Path, r.Version)
}

type HTTPResponse struct {
	Status        HTTPStatusCode
	Server        string
	Date          string
	Content       []byte
	ContentType   string
	ContentLength int
}

func (r HTTPResponse) Bytes() []byte {
	var response_buffer []byte = []byte(fmt.Sprintf("HTTP/1.0 %s\nServer: %s\nDate: %s\nContent-type: %s\nContent-Length: %d\n\n",
		r.Status,
		r.Server,
		r.Date,
		r.ContentType,
		r.ContentLength))

	return append(response_buffer, r.Content...)
}

var httpMethodHandlers map[string]func(HTTPRequest) (HTTPResponse, error) = map[string]func(HTTPRequest) (HTTPResponse, error){"GET": handleGETRequest}

var contentTypes map[string]string = map[string]string{
	"txt":  "text/plain",
	"html": "text/html",
	"css":  "text/css",
	"js":   "text/javascript",

	"jpeg": "image/jpeg",
	"jpg":  "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",

	"mpeg": "audio/mpeg",
	"mp4":  "video/mp4",

	"json": "application/json",
	"xml":  "application/xml",
	"pdf":  "application/pdf",
}

func parseStartLine(start_line string) (string, string, string) {
	split_result := strings.Split(start_line, " ")
	http_method, path, version := split_result[0], split_result[1], strings.TrimSuffix(split_result[2], "\r")
	return http_method, path, version
}

func parseHTTPRequest(raw_request []byte) (HTTPRequest, error) {
	var request string = string(raw_request)
	var http_request HTTPRequest
	var err error

	split_request := strings.Split(request, "\n")

	http_request.Method, http_request.Path, http_request.Version = parseStartLine(split_request[0])

	return http_request, err
}
func generateHTTPResponse(status HTTPStatusCode, content []byte, content_type string) HTTPResponse {
	return HTTPResponse{
		Status:        status,
		Server:        "Eldad's GO HTTP Server",
		Date:          time.Now().String(),
		Content:       content,
		ContentType:   content_type,
		ContentLength: len(content),
	}
}
func handleGETRequest(http_request HTTPRequest) (HTTPResponse, error) {
	var http_response HTTPResponse
	var file_path string = http_request.Path

	if file_path == "/" {
		file_path = "/index.html"
	}

	if _, err := os.Stat(file_path[1:]); errors.Is(err, os.ErrNotExist) {
		log.Printf("%s %s\n", http_request, HTTP_NOT_FOUND)
		http_response = generateHTTPResponse(HTTP_NOT_FOUND, default404Page, "text/html; charset=utf-8")
		return http_response, nil
	}

	file_content, err := os.ReadFile(file_path[1:])
	if err != nil {
		fmt.Println(err)
		return http_response, err
	}
	file_extension := path.Ext(file_path[1:])
	content_type, ok := contentTypes[file_extension[1:]]
	if !ok {
		content_type = "application/octet-stream"
	}
	log.Printf("%s %s\n", http_request, HTTP_OK)
	http_response = generateHTTPResponse(HTTP_OK, file_content, content_type)

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
		http_response, err = handler(http_request)
	} else {
		log.Printf("Got a %s Request, which is not supported\n", http_request.Method)
		err = fmt.Errorf("unsupported request: %s", http_request.Method)
	}
	return http_response, err

}
