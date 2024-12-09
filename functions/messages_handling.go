package netcat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

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
