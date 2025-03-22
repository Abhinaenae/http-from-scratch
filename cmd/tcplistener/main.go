package main

import (
	"fmt"
	"log"
	"net"

	"github.com/abhinaenae/http-from-scratch/internal/request"
)

const netPort = ":8000"

func main() {
	listener, err := net.Listen("tcp", netPort)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()
	fmt.Println("Listening for TCP traffic on", netPort)
	for {
		// Wait for a connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Connection established with %s\n", conn.RemoteAddr())
		fmt.Println("=====================================")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error parsing request: %s\n", err.Error())
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Printf("Connection to %s closed\n", conn.RemoteAddr())
	}

}
