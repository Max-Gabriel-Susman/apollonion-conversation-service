package main

import (
	"net"
	"testing"
	"time"
)

func TestChatServer(t *testing.T) {
	go main()               // Start the server in a goroutine
	time.Sleep(time.Second) // Give the server a moment to start

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("Failed to connect to chat server: %s", err)
	}
	defer conn.Close()

	message := "Hello, World!\n"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to send message to chat server: %s", err)
	}

	// Here you would extend the test to read back the message from the server
	// This is left as an exercise for the reader due to complexity
}
