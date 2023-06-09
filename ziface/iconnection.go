package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetConnId() uint32
	GetTcpConnection() *net.TCPConn
	RemoteAddr() net.Addr
}

type HandleFunc func(*net.TCPConn, []byte, int) error
