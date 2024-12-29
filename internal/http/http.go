package http

import (
	_ "embed" // Used to embed the contents of the default404Page file into a variable
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type httpStatusCode string

const (
	httpOk       httpStatusCode = "200 OK"
	httpNotFound httpStatusCode = "404 File not found"
)

//go:embed templates/page_not_found.html
var default404Page []byte

// Request represents the data of an HTTP request
type Request struct {
	Method  string
	Path    string
	Version string
}

func (r Request) String() string {
	return fmt.Sprintf("%s %s %s", r.Method, r.Path, r.Version)
}

// Response represents data of an HTTP response
type Response struct {
	Status        httpStatusCode
	Server        string
	Date          string
	Content       []byte
	ContentType   string
	ContentLength int
}

// Bytes structures data from the Response struct to a bytes slice
func (r Response) Bytes() []byte {
	var responseBuffer []byte = []byte(fmt.Sprintf("HTTP/1.0 %s\nServer: %s\nDate: %s\nContent-type: %s\nContent-Length: %d\n\n",
		r.Status,
		r.Server,
		r.Date,
		r.ContentType,
		r.ContentLength))

	return append(responseBuffer, r.Content...)
}

var httpMethodHandlers map[string]func(Request) (Response, error) = map[string]func(Request) (Response, error){"GET": handleGETRequest}

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

func parseStartLine(startLine string) (string, string, string) {
	splitResult := strings.Split(startLine, " ")
	httpMethod, path, version := splitResult[0], splitResult[1], strings.TrimSuffix(splitResult[2], "\r")
	return httpMethod, path, version
}

func parseHTTPRequest(rawRequest []byte) (Request, error) {
	var request string = string(rawRequest)
	var httpRequest Request
	var err error

	splitRequest := strings.Split(request, "\n")

	httpRequest.Method, httpRequest.Path, httpRequest.Version = parseStartLine(splitRequest[0])

	return httpRequest, err
}
func generateHTTPResponse(status httpStatusCode, content []byte, contentType string) Response {
	return Response{
		Status:        status,
		Server:        "Eldad's GO HTTP Server",
		Date:          time.Now().String(),
		Content:       content,
		ContentType:   contentType,
		ContentLength: len(content),
	}
}
func handleGETRequest(httpRequest Request) (Response, error) {
	var httpResponse Response
	var filePath string = httpRequest.Path

	if filePath == "/" {
		filePath = "/index.html"
	}

	if _, err := os.Stat(filePath[1:]); errors.Is(err, os.ErrNotExist) {
		log.Printf("%s %s\n", httpRequest, httpNotFound)
		httpResponse = generateHTTPResponse(httpNotFound, default404Page, "text/html; charset=utf-8")
		return httpResponse, nil
	}

	fileContent, err := os.ReadFile(filePath[1:])
	if err != nil {
		fmt.Println(err)
		return httpResponse, err
	}
	fileExtension := path.Ext(filePath[1:])
	contentType, ok := contentTypes[fileExtension[1:]]
	if !ok {
		contentType = "application/octet-stream"
	}
	log.Printf("%s %s\n", httpRequest, httpOk)
	httpResponse = generateHTTPResponse(httpOk, fileContent, contentType)

	return httpResponse, err
}

// HandleRequest handles a raw http request and returns a valid Resposne
func HandleRequest(rawRequest []byte) (Response, error) {
	var httpResponse Response
	httpRequest, err := parseHTTPRequest(rawRequest)
	if err != nil {
		return httpResponse, err
	}

	handler, ok := httpMethodHandlers[httpRequest.Method]

	if ok {
		httpResponse, err = handler(httpRequest)
	} else {
		log.Printf("Got a %s Request, which is not supported\n", httpRequest.Method)
		err = fmt.Errorf("unsupported request: %s", httpRequest.Method)
	}
	return httpResponse, err

}
