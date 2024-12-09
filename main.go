package main

import n "netcat/functions"

func main() {
	listener := n.StartServer()
	if listener == nil {
		return
	}
	defer listener.Close()

	n.AcceptConnections(listener)
}
