package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go processRequest(conn)
	}

}

func processRequest(conn net.Conn) {

	defer conn.Close()

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}
	req := string(buff[:n])

	if strings.HasPrefix(req, "GET / HTTP/1.1") {
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}

	} else if strings.HasPrefix(req, "GET /files/") {
		reqLines := strings.Split(req, "\r\n")
		target := strings.Split(reqLines[0], " ")[1]
		filename := strings.Split(target, "/")[2]
		path := "/tmp/data/codecrafters.io/http-server-tester/" + filename
		var message string
		// Read the entire file
		fileInfo, err := os.Stat(path)
		if err != nil {
			message = "HTTP/1.1 404 Not Found\r\n\r\n"
		} else {
			data, _ := os.ReadFile(path)
			message = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %v\r\n\r\n%v", fileInfo.Size(), string(data))
		}
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}

	} else if strings.HasPrefix(req, "GET /user-agent") {

		reqLines := strings.Split(req, "\r\n")
		var userAgent string
		for _, v := range reqLines {
			if strings.HasPrefix(v, "User-Agent") {
				userAgent = strings.Split(v, " ")[1]
			}
		}

		message := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(userAgent), userAgent)

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(req, "GET /echo/") {
		reqLines := strings.Split(req, "\r\n")
		target := strings.Split(reqLines[0], " ")[1]
		responseBody := strings.Split(target, "/")[2]
		message := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(responseBody), responseBody)
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(req, "POST") {
		reqLines := strings.Split(req, "\r\n")

		body := reqLines[5]
		filename := strings.Split(strings.Split(reqLines[0], " ")[1], "/")[2]
		path := "/tmp/data/codecrafters.io/http-server-tester/" + filename
		f, _ := os.Create(path)
		defer f.Close()
		f.WriteString(body)

		conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))

	} else {
		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	}
}
