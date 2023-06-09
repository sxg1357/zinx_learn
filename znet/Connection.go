package znet

import (
	"errors"
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
	//HandleFunc     ziface.HandleFunc
	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connId uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:           conn,
		ConnId:         connId,
		IsClose:        false,
		ExitBufferChan: make(chan bool),
		Router:         router,
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
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		cnt, err := io.ReadFull(c.GetTcpConnection(), headData)
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
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error")
			c.ExitBufferChan <- true
			return
		}
		var data = make([]byte, msg.GetDataLen())
		if msg.GetDataLen() > 0 {
			_, err := io.ReadFull(c.GetTcpConnection(), data)
			if err != nil {
				c.ExitBufferChan <- true
				continue
			}
		}
		msg.SetData(data)
		req := Request{
			conn: c,
			data: msg,
		}
		go func(request ziface.IRequest) {
			c.Router.PreHandler(request)
			c.Router.Handler(request)
			c.Router.PostHandler(request)
		}(&req)
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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClose == true {
		return errors.New("client has already closed")
	}
	dp := NewDataPack()
	dataPack, err := dp.Pack(NewMessage(data, msgId))
	if err != nil {
		return err
	}
	_, err = c.Conn.Write(dataPack)
	if err != nil {
		c.ExitBufferChan <- true
		return err
	}
	return nil
}
