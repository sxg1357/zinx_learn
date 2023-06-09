package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetConnId() uint32
}

type HandleFunc func(*net.TCPConn, []byte, int) error
