package netcat

import (
	"fmt"
	"net"
	"os"
)

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
		fmt.Println("Error listening:", err)
		return nil
	}
	addr := listener.Addr().String()
	_, port, err = net.SplitHostPort(addr)
	if err != nil {
		fmt.Printf("error spliting the host from the port: %v\n", err)
		return nil
	}
	fmt.Printf("Server started on port: %s\n", port)

	return listener
}

func AcceptConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connections:", err)
			continue
		}
		go HandleClients(conn)
	}
}
