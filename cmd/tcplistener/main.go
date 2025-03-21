package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
		linesChan := getLinesChannel(conn)
		for line := range linesChan {
			fmt.Println(line)
		}
		fmt.Printf("Connection to %s closed\n", conn.RemoteAddr())
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {

	ch := make(chan string)
	go func() {
		defer f.Close()
		defer close(ch)
		currentLineContents := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					ch <- currentLineContents // Send last line if file doesn't end with newline
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := range len(parts) - 1 {
				line := currentLineContents + parts[i]
				ch <- line
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}

		close(ch)
	}()

	return ch
}
