package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	srvAddr := "localhost:8000"
	addr, err := net.ResolveUDPAddr("udp", srvAddr)
	if err != nil {
		log.Fatalf("Could not resolve udp address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Could not dial udp address: %v", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading line: %v\n", err)
			continue
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Printf("Error writing line: %v\n", err)
		}
	}
}
