package main

import (
	"bufio"
	"net"
	"strings"
	"testing"
)

func TestChatServer(t *testing.T) {
	go main()

	t.Run("Successful message broadcast", func(t *testing.T) {
		conn, err := net.Dial("tcp", "localhost:8081")
		if err != nil {
			t.Fatalf("Failed to connect to chat server: %s", err)
		}
		defer conn.Close()

		message := "Hello, World!\n"
		_, err = conn.Write([]byte(message))
		if err != nil {
			t.Fatalf("Failed to send message to chat server: %s", err)
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read response from chat server: %s", err)
		}

		expectedSuffix := strings.TrimSpace(message)
		if !strings.HasSuffix(strings.TrimSpace(response), expectedSuffix) {
			t.Fatalf("Expected server response to end with %q, got %q", expectedSuffix, response)
		}
	})

	t.Run("Failure on unexpected message format", func(t *testing.T) {
		conn, err := net.Dial("tcp", "localhost:8081")
		if err != nil {
			t.Fatalf("Failed to connect to chat server: %s", err)
		}
		defer conn.Close()

		message := "InvalidMessageFormat\n"
		_, err = conn.Write([]byte(message))
		if err != nil {
			t.Fatalf("Failed to send message to chat server: %s", err)
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read response from chat server: %s", err)
		}

		if !strings.Contains(response, strings.TrimSpace(message)) {
			t.Fatalf("Expected server response to contain %q, got %q", message, response)
		}
	})

}
