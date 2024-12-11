package main

import (
	n "netcat/functions"
)

func main() {
	n.ClearFile("file/messages.txt")
	listener := n.StartServer()
	if listener == nil {
		return
	}
	n.AcceptConnections(listener)
}
