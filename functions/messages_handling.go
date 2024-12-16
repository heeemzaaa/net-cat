package netcat

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// this function handles the messages sent by the client 
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
			_, err := conn.Write([]byte("Invalid message !\n"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error sending the invalid message: %v\n", err)
				return
			}
		}
		if message != "" && len(message) < 300 {
			timeNow := time.Now().Format("2006-01-02 15:04:05")
			formattedMessage := fmt.Sprintf("\n[%s][%s]:%s\n", timeNow, name, message)
			BroadcastMessage(formattedMessage, conn)
		} else {
			TypingPlace(name, conn)
		}

	}
}

// this function handles the prompt giving to the user  
func TypingPlace(name string, conn net.Conn) {
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	prompt := fmt.Sprintf("[%s][%s]:", timeNow, name)
	_, err := conn.Write([]byte(prompt))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending typing prompt to %s: %v\n", name, err)
	}
}


