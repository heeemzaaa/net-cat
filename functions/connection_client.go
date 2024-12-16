package netcat

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// this function handles the client connection

func HandleClients(conn net.Conn) {
	defer conn.Close()

	Welcoming(conn)
	name := ClientName(conn)
	if counter >= 10 {
		RejectConnection(conn)
		return
	}
	IncrementConnectionCount()
	SendHistoryChat(conn)
	message := fmt.Sprintf("\n%s has joined our chat...\n", name)
	BroadcastMessage(message, conn)
	Handlemessages(name, conn)
}

// this function welcomes the new client by the logo and the prompt
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
		fmt.Fprintln(os.Stderr, "Error sending the welcoming message:", err)
		return
	}
}

// this function handles the client name
func ClientName(conn net.Conn) string {
	logCount := 0
	for {
		if logCount > 5 {
			_, err := conn.Write([]byte("\nYou've reached your trying limit"))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error sending the message:", err)
				return ""
			}
			logCount = 0
			conn.Close()
		}
		reader := bufio.NewReader(conn)
		name, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading the name from the user:", err)
			return ""
		}

		name = strings.TrimSpace(name)
		if !Printable(name) {
			name = ""
		}

		if name == "" || len(name) > 25 {
			conn.Write([]byte("Please enter a valid name!\nExample of a valid name : oumayma , yassine , hamza\n[ENTER YOUR NAME]: "))
			logCount++
			continue
		}

		if SpaceName(name) {
			conn.Write([]byte("The username must be without space !\nExample of a valid name : oumayma , yassine , hamza\n[ENTER YOUR NAME]: "))
			logCount++
			continue
		}

		mu.Lock()
		if _, exists := clientM[name]; exists {
			mu.Unlock()
			conn.Write([]byte("Name already taken. Choose another one:\nExample of a valid name : oumayma , yassine , hamza\n[ENTER YOUR NAME]: "))
			logCount++
			continue
		}
		clientM[name] = conn
		mu.Unlock()

		return name
	}
}

// this function broadcast the message to other clients
func BroadcastMessage(message string, excluded net.Conn) {
	mu.Lock()
	StoreMessages(message)
	defer mu.Unlock()

	for clientName, conn := range clientM {
		if conn == excluded {
			TypingPlace(clientName, conn)
			continue
		} else {
			_, err := conn.Write([]byte(message))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error broadcasting to %s: %v\n", clientName, err)
				continue
			}
			TypingPlace(clientName, conn)
		}
	}
}

// this function checks if the client isn't there and removes it
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

// this function increment the connection
func IncrementConnectionCount() {
	mu.Lock()
	defer mu.Unlock()
	counter++
}

// this function decrement the connection
func DecrementConnectionCount() {
	mu.Lock()
	defer mu.Unlock()
	counter--
}

// this function reject connection if the connection reaches its limit
func RejectConnection(conn net.Conn) {
	_, err := conn.Write([]byte("The room has reached its limit , try again later !\n"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error sending the rejection message:", err)
		return
	}
	conn.Close()
}
