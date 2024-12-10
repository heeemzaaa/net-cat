package netcat

import (
	"net"
	"sync"
)

var (
	clientM        = make(map[string]net.Conn)
	mu             sync.Mutex
	counter        int
)
