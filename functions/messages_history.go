package netcat

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// this function clears the file each time the program starts
func ClearFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error clearing file %s: %v\n", filename, err)
		return
	}
	defer file.Close()
}

// this function stores the messages in the file
func StoreMessages(message string) {
	filePath := "file/messages.txt"
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening the file:", err)
		return
	}
	defer file.Close()

	message = strings.TrimSpace(message)
	_, err = file.WriteString(message + "\n")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing in the file :", err)
		return
	}
}

// this function sends the history chat
func SendHistoryChat(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	filePath := "file/messages.txt"
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from the file:", err)
		return
	}

	storedMessages := strings.Split(string(content), "\n")

	for _, message := range storedMessages {
		if len(message) > 0 && message[0] == '[' {
			slice := strings.Split(message, ":")

			if len(slice) > 3 && slice[3] != "" {
				_, err := conn.Write([]byte(message + "\n"))
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error sending chat history: %v\n", err)
					return
				}
			}
		}
	}
}
