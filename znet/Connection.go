package znet

import (
	"fmt"
	"io"
	"net"
	"zinx_learn/ziface"
)

type Connection struct {
	Conn           *net.TCPConn
	ConnId         uint32
	IsClose        bool
	ExitBufferChan chan bool
	HandleFunc     ziface.HandleFunc
}

func NewConnection(conn *net.TCPConn, connId uint32, HandlerApi ziface.HandleFunc) *Connection {
	return &Connection{
		Conn:           conn,
		ConnId:         connId,
		IsClose:        false,
		ExitBufferChan: make(chan bool),
		HandleFunc:     HandlerApi,
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	for {
		select {
		case <-c.ExitBufferChan:
			fmt.Println("Reader Goroutine exit")
			return
		}
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if cnt == 0 {
			fmt.Println("client close...")
			c.ExitBufferChan <- true
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("read error")
			c.ExitBufferChan <- true
			continue
		}
		if err := c.HandleFunc(c.Conn, buf, cnt); err != nil {
			fmt.Println("call Handler Api error...")
			c.ExitBufferChan <- true
			return
		}
	}
}

func (c *Connection) Stop() {
	fmt.Printf("ConnId:%d stop\r\n", c.ConnId)
	if c.IsClose == true {
		return
	}
	c.IsClose = true
	c.ExitBufferChan <- true
	c.Conn.Close()
}

func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
