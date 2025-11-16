package parser

import (
	"errors"
	//"fmt"
	"slices"
	//"strconv"
	"strings"
	//"time"
)

func ParseRequestHeaders(lines <-chan string) (*RequestHeaders, error) {
	firstLine := <- lines
	method, path, version, err := parseFirstLine(firstLine)
	if err != nil {
		return nil, err
	}
	headers := map[string]string{}
	for line := range lines {
		//fmt.Println("read: ", line)
		if(line == "" || line == "\r") {
			break
		}
		key, value, err := parseHeader(line)
		if err != nil {
			continue
		}
		headers[key] = strings.TrimSpace(value)
	}

	request := &RequestHeaders{
		Method: method,
		Path: path,
		HttpVersion: version,
		Headers: headers,
	}
    /*

	val, ok := request.Headers["Content-Length"]
	if ok {

		contentLength, err := strconv.Atoi(val)
		if err {
			return request, err
		}

		body := readBody(lines, contentLength)
	}
	*/
	return request, nil
}

func parseHeader(line string) (string, string, error) {
	splitLine := strings.SplitN(line, ":", 2)
	if len(splitLine) < 2 {
		return "", "", errors.New("Header formatted incorrectly")
	}
	return splitLine[0], splitLine[1], nil
}

func parseFirstLine(firstLine string) (string, string, string, error) {
	splitLine := strings.Split(firstLine, " ")
	if len(splitLine) != 3 {
		return "", "", "", errors.New("First line formatted incorrectly")
	}

	method, path, version := strings.ToUpper(splitLine[0]), splitLine[1], splitLine[2]

	if !isValidHttpMethod(method) {
		return "", "", "", errors.New("Invalid HTTP method")
	}

	if !isValidPath(path) {
		return "", "", "", errors.New("Invalid path")
	}

	if !isValidHttpVersion(version) {
		return "", "", "", errors.New("Unsupported HTTP version")
	}

	return method, path, version, nil
}



func isValidHttpMethod(method string) bool {
	VALID_HTTP_METHODS := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
	return slices.Contains(VALID_HTTP_METHODS, method)
}

func isValidPath(path string) bool {
	// TODO: Implement regex
	return true
}

func isValidHttpVersion(version string) bool {
	// TODO: Implement
	return true
}
