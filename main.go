package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	clientM        = make(map[string]net.Conn)
	storedMessages = []string{}
	mu             sync.Mutex
	counter        int
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
	fmt.Printf("Server started on port %s\n", port)
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

func HandleClients(conn net.Conn) {
	defer conn.Close()

	Welcoming(conn)
	name := ClientName(conn)
	IncrementConnectionCount()
	if counter >= 3 {
		RejectConnection(conn)
		return
	}
	SendHistoryChat(conn)
	message := fmt.Sprintf("\n%s has joined our chat...\n", name)
	BroadcastMessage(message, conn)
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
		fmt.Println("Error sending the welcoming message:", err)
		return
	}
}

func ClientName(conn net.Conn) string {
	logCount := 0
	for {
		if logCount > 2 {
			_, err := conn.Write([]byte("\nYou've reached your trying limit"))
			if err != nil {
				fmt.Println("Error:", err)
				return ""
			}
			logCount = 0
			conn.Close()
		}
		reader := bufio.NewReader(conn)
		name, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading the name from the user:", err)
			return ""
		}

		name = strings.TrimSpace(name)
		if !Printable(name) {
			name = ""
		}

		if name == "" || len(name) > 25 {
			conn.Write([]byte("Please enter a valid name!\n[ENTER YOUR NAME]: "))
			logCount++
			continue
		}

		if SpaceName(name) {
			conn.Write([]byte("The username must be without space !\n[ENTER YOUR NAME]: "))
			logCount++
			continue
		}

		mu.Lock()
		if _, exists := clientM[name]; exists {
			mu.Unlock()
			conn.Write([]byte("Name already taken. Choose another one:\n[ENTER YOUR NAME]: "))
			logCount++
			continue
		}
		clientM[name] = conn
		mu.Unlock()

		fmt.Printf("Client %s connected\n", name)
		return name
	}
}

func BroadcastMessage(message string, excluded net.Conn) {
	mu.Lock()
	storedMessages = append(storedMessages, strings.TrimSpace(message))
	defer mu.Unlock()

	for clientName, conn := range clientM {
		if conn == excluded {
			TypingPlace(clientName, conn)
			continue
		} else {
			_, err := conn.Write([]byte(message))
			if err != nil {
				fmt.Printf("Error broadcasting to %s: %v\n", clientName, err)
				continue
			}
			TypingPlace(clientName, conn)
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

	message := fmt.Sprintf("\n%s has left our chat...\n", name)
	BroadcastMessage(message, nil)
}

func Handlemessages(name string, conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Client %s disconnected: %v\n", name, err)
			RemoveClient(name)
			DecrementConnectionCount()
			return
		}

		message = strings.TrimSpace(message)
		if !Printable(message) {
			message = ""
		}
		if message != "" && len(message) < 300 {
			timeNow := time.Now().Format("2006-01-02 15:04:05")
			formattedMessage := fmt.Sprintf("\n[%s][%s]: %s\n", timeNow, name, message)
			BroadcastMessage(formattedMessage, conn)
		} else {
			TypingPlace(name, conn)
		}

	}
}

func TypingPlace(name string, conn net.Conn) {
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	prompt := fmt.Sprintf("[%s][%s]:", timeNow, name)
	_, err := conn.Write([]byte(prompt))
	if err != nil {
		fmt.Printf("Error sending typing prompt to %s: %v\n", name, err)
	}
}

func StoreMessages(message string) {
	mu.Lock()
	storedMessages = append(storedMessages, message)
	mu.Unlock()
}

func SendHistoryChat(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	for i := 0; i < len(storedMessages); i++ {
		if storedMessages[i][0] == '[' && len(storedMessages[i]) != 0 {
			slice := strings.Split(storedMessages[i], ":")
			if slice[3] != "" {
				_, err := conn.Write([]byte(storedMessages[i] + "\n"))
				if err != nil {
					fmt.Printf("Error sending chat history: %v\n", err)
					return
				}
			}

		}
	}
}

func IncrementConnectionCount() {
	mu.Lock()
	defer mu.Unlock()
	counter++
}

func DecrementConnectionCount() {
	mu.Lock()
	defer mu.Unlock()
	counter--
}

func RejectConnection(conn net.Conn) {
	_, err := conn.Write([]byte("The room has reached its limit , try again later !\n"))
	if err != nil {
		fmt.Println("Error :", err)
		return
	}
	conn.Close()
}

func Printable(message string) bool {
	slice := []rune(message)
	for i := 0; i < len(slice); i++ {
		if !(slice[i] >= 32 && slice[i] <= 126) {
			return false
		}
	}
	return true
}

func SpaceName(name string) bool {
	slice := []rune(name)
	for i := 0; i < len(slice); i++ {
		if slice[i] == ' ' {
			return true
		}
	}
	return false
}

func main() {
	listener := StartServer()
	if listener == nil {
		return
	}
	defer listener.Close()

	// Graceful shutdown handling
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		fmt.Println("\nShutting down the server...")
		mu.Lock()
		for _, conn := range clientM {
			conn.Close()
		}
		mu.Unlock()
		listener.Close()
		os.Exit(0)
	}()

	AcceptConnections(listener)
}
