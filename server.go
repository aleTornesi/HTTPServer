package main

import (
	"HTTP/internal/parser"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
)

func main() {
	Serve()
}

func Serve() {
	listener, err := net.Listen("tcp", ":8080")


	if err != nil {
		//error handling
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			//error handling
			continue
		}
		go handleConnection(conn)
	}
}


func handleConnection(listener net.Conn) {
	lines := getLinesFromReader(listener)
	requestHeaders, err := parser.ParseRequestHeaders(lines)
	if err != nil {
		return
	}

	val, ok := requestHeaders.Headers["Content-Length"]
	var body string
	if ok {
		contentLength, err := strconv.Atoi(val)
		if err == nil{
			body = readBody(listener, contentLength)
		}
	}
	fmt.Println("Body: ", body)
}

func readBody(source net.Conn, contentLength int) string {
	defer source.Close()
	s := ""
	for contentLength > 0 {
		chunck := make([]byte, min(contentLength, 8))
		size, err := source.Read(chunck)
		if err != nil {
			break
		}

		chunck = chunck[:size]
		contentLength -= size
		s += string(chunck)
	}
	return s
}

func getLinesFromReader(r io.ReadCloser) <- chan string {
	out := make(chan string, 1)

	go func() {
		defer close(out)

		s := ""
		foundSeparator := false
		for !foundSeparator {
			chunck := make([]byte, 8)
			size, err := r.Read(chunck)
			if err != nil {
				r.Close()
				break
			}
			chunck = chunck[:size]
			newLineIndex := bytes.IndexByte(chunck, '\n')
			for newLineIndex != -1 {
				s += string(chunck[:newLineIndex])
				chunck = chunck[newLineIndex + 1:]
				if s == "" || s == "\r" {
					foundSeparator = true
				}
				out <- s
				s = ""
				newLineIndex = bytes.IndexByte(chunck, '\n')
			}

			s += string(chunck)
		}
		if len(s) != 0 {
			out <- s
		}
	}()

	return out
}

