package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abhinaenae/http-from-scratch/internal/server"
)

const port = 8000

func main() {
	server, err := server.Serve(port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
