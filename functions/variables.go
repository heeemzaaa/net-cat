package netcat

import (
	"net"
	"sync"
)

var (
	clientM        = make(map[string]net.Conn)
	storedMessages = []string{}
	mu             sync.Mutex
	counter        int
)
