package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Fprintf(conn, "Received: %s\n", text)
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

/*

	set the stage for channels:
		store identifiers and broadcast messages to all clients

		the client that sends the message should not receive it back

		test coverage

		should open up the door for using channels and then we can

		bring back in a good implementation of fanout orchestration

	controlling concurrency:
		cap goroutines, but still accept more clients to connect

		maybe store a slice of net.Conn to sorta cache them so we can get

		around to them later while staying under the goroutine cap

		test coverage

	failure case coverage for existing and new logic:
*/
