package netcat

import (
	"fmt"
	"net"
	"os"
)

// this function starts a server and gives the port to the user
func StartServer() net.Listener {
	port := ":8989"
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return nil
	}
	if len(os.Args) == 2 {
		port = ":" + os.Args[1]
	}
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error listening:", err)
		return nil
	}
	addr := listener.Addr().String()
	_, port, err = net.SplitHostPort(addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error spliting the host from the port: %v\n", err)
		return nil
	}
	fmt.Printf("Server started on port: %s\n", port)

	return listener
}

// this function accept the connections
func AcceptConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error accepting connections:", err)
			continue
		}
		go HandleClients(conn)
	}
}
