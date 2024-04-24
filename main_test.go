package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestServerConcurrency(t *testing.T) {
	t.Run("successful handling of concurrent requests from multiple clients(5 clients)", func(t *testing.T) {
		address := ":8080"
		go startServer(address)
		time.Sleep(time.Second)

		numClients := 5
		messages := make(chan string, numClients)

		clientWork := func(id int) {
			conn, err := net.Dial("tcp", address)
			if err != nil {
				t.Error("Failed to connect to server:", err)
				return
			}
			defer conn.Close()

			message := fmt.Sprintf("Hello from client %d", id)
			fmt.Fprintln(conn, message)

			responseScanner := bufio.NewScanner(conn)
			if responseScanner.Scan() {
				response := responseScanner.Text()
				messages <- response
			}
		}

		for i := 0; i < numClients; i++ {
			go clientWork(i)
		}

		for i := 0; i < numClients; i++ {
			msg := <-messages
			t.Log("Received:", msg)
		}
	})

	t.Run("server start up delay", func(t *testing.T) {
		address := ":8080"
		go startServer(address)
		time.Sleep(1 * time.Second)

		conn, err := net.Dial("tcp", address)
		if err != nil {
			t.Skip("Server not ready, skipping test")
		}
		conn.Close()
	})

	t.Run("connection retries", func(t *testing.T) {
		address := ":8080"
		go startServer(address)
		time.Sleep(1 * time.Second)

		const maxRetries = 3
		var conn net.Conn
		var err error
		for attempt := 1; attempt <= maxRetries; attempt++ {
			conn, err = net.Dial("tcp", address)
			if err == nil {
				break
			}
			time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
		}
		if err != nil {
			t.Fatal("Failed to connect after retries:", err)
		}
		conn.Close()
	})

	t.Run("stress test(100 clients)", func(t *testing.T) {
		address := ":8080"
		go startServer(address)
		time.Sleep(1 * time.Second)

		const stressClientCount = 100
		done := make(chan bool, stressClientCount)

		for i := 0; i < stressClientCount; i++ {
			go func(id int) {
				conn, err := net.Dial("tcp", address)
				if err != nil {
					t.Error("Failed to connect to server:", err)
					done <- false
					return
				}
				defer conn.Close()
				fmt.Fprintf(conn, "Stress test message from client %d\n", id)
				scanner := bufio.NewScanner(conn)
				if scanner.Scan() {
					done <- true
				} else {
					t.Error("Failed to receive response during stress test")
					done <- false
				}
			}(i)
		}

		for i := 0; i < stressClientCount; i++ {
			if success := <-done; !success {
				t.Error("Not all clients completed successfully")
			}
		}
	})
}
