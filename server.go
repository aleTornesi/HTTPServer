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

	fmt.Println("HEADERS", requestHeaders.Headers)
	val, ok := requestHeaders.Headers["Content-Length"]
	fmt.Println(ok, val)
	var body string
	if ok {
		contentLength, err := strconv.Atoi(val)
		if err != nil{
			fmt.Printf("Couldnt parse Content-Length (%s)", val)
		}
		body = readBody(listener, contentLength)
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
			fmt.Println(err)
		}

		chunck = chunck[:size]
		fmt.Println(string(chunck))
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
			//fmt.Println("Read chunck", string(chunck), chunck)
			newLineIndex := bytes.IndexByte(chunck, '\n')
			for newLineIndex != -1 {
				s += string(chunck[:newLineIndex])
				chunck = chunck[newLineIndex + 1:]
				//fmt.Println("Remaining chunck:", string(chunck), chunck, "inserting ", s, "in the channel")
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

