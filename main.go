package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	clientM = make(map[string]net.Conn)
	mu      sync.Mutex
)

func StartServer() net.Listener {
	port := ":8989"
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port\n")
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
	return listener
}

func AcceptConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connections:", err)
			return
		}
		go HandleClients(conn)
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleClients(conn net.Conn) {
	defer conn.Close()

	Welcoming(conn)
	name := ClientName(conn)
	Handlemessages(name, conn)
}

func Welcoming(conn net.Conn) {
	welcome := `Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    '.       | '' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     '-'       '--'
[ENTER YOUR NAME]: `

	_, err := conn.Write([]byte(welcome))
	if err != nil {
		fmt.Println("Error sending the welcoming message :", err)
		return
	}
}

func ClientName(conn net.Conn) string {
	for {
		reader := bufio.NewReader(conn)
		name, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading the name from the user:", err)
			return ""
		}

		name = strings.TrimSpace(name)

		if name == "" || len(name) > 25 {
			conn.Write([]byte("Please enter a valid name!\n"))
			continue
		}

		mu.Lock()
		if _, exists := clientM[name]; exists {
			mu.Unlock()
			conn.Write([]byte("Name already taken. Choose another one:\n"))
			continue
		}

		clientM[name] = conn
		mu.Unlock()
		message := fmt.Sprintf("\n%s has joined our chat...\n", name)
		BroadcastMessage(message, conn)
		fmt.Printf("Client %s connected\n", name)
		return name
	}
}

func BroadcastMessage(message string, exluded net.Conn) {
	for client, conn := range clientM {
		if conn != exluded {
			_, err := conn.Write([]byte(message))
			if err != nil {
				fmt.Printf("Error broadcasting the message to the client %s : %v\n", client, err)
				return
			}
		}
	}
}

func RemoveClient(name string) {
	mu.Lock()
	conn, ok := clientM[name]
	if !ok {
		mu.Unlock()
		fmt.Printf("Client %s not found.\n", name)
		return
	}

	conn.Close()
	delete(clientM, name)
	mu.Unlock()

	message := fmt.Sprintf("%s has left our chat...\n", name)
	BroadcastMessage(message, nil)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func Handlemessages(name string, conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Client %s disconnected: %v\n", name, err)
			RemoveClient(name)
			return
		}

		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}
		timeNow := time.Now().Format(time.DateTime)
		formattedMessage := fmt.Sprintf("[%s][%s]:", timeNow, name)
		BroadcastMessage(formattedMessage, conn)
	}
}

func IsEmpty(message string) bool {
	if message == "" {
		return true
	}
	return false
}

func main() {
	listener := StartServer()
	if listener == nil {
		return
	}
	defer listener.Close()

	AcceptConnections(listener)
}
