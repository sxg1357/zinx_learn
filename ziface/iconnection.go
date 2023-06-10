package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetConnId() uint32
	GetTcpConnection() *net.TCPConn
	RemoteAddr() net.Addr
	SendMsg(msgId uint32, data []byte) error
}

//type HandleFunc func(*net.TCPConn, []byte, int) error
