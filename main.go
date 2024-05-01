package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

var (
	connMutex     sync.Mutex
	connections   []net.Conn
	broadcastChan = make(chan string)
)

func broadcaster() {
	for {
		msg := <-broadcastChan
		connMutex.Lock()
		for _, conn := range connections {
			fmt.Fprintf(conn, "Broadcast: %s\n", msg)
		}
		connMutex.Unlock()
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		connMutex.Lock()
		for i, c := range connections {
			if c == conn {
				connections = append(connections[:i], connections[i+1:]...)
				break
			}
		}
		connMutex.Unlock()
	}()

	connMutex.Lock()
	connections = append(connections, conn)
	connMutex.Unlock()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		broadcastChan <- text
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading from connection: %s\n", err)
	}
}

func startServer(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Println("Server is listening on", address)

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error accepting connection: %s\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func main() {
	if err := startServer(":8080"); err != nil {
		fmt.Fprintf(os.Stderr, "error starting server: %s\n", err)
		os.Exit(1)
	}
}
