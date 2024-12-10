package netcat

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
	filePath := "file/messages.txt"

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0o666)
	if err != nil {
		fmt.Println("Error opening the file :", err)
		return
	}
	defer file.Close()
	message = strings.TrimSpace(message)
	_, err = file.WriteString(message + "\n")
	if err != nil {
		fmt.Println("Error writing in the file :", err)
		return
	}
}

func SendHistoryChat(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	filePath := "file/messages.txt"
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading from the file:", err)
		return
	}

	storedMessages := strings.Split(string(content), "\n")

	for _, message := range storedMessages {
		if len(message) > 0 && message[0] == '[' {
			slice := strings.Split(message, ":")

			if len(slice) > 3 && slice[3] != "" {
				_, err := conn.Write([]byte(message + "\n"))
				if err != nil {
					fmt.Printf("Error sending chat history: %v\n", err)
					return
				}
			}
		}
	}
}
