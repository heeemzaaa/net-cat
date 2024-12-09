package netcat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func HandleClients(conn net.Conn) {
	defer conn.Close()

	Welcoming(conn)
	name := ClientName(conn)
	IncrementConnectionCount()
	if counter >= 10 {
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
