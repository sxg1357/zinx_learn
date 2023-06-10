package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx_learn/utils"
	"zinx_learn/ziface"
)

type Connection struct {
	Conn           *net.TCPConn
	ConnId         uint32
	IsClose        bool
	MsgHandler     ziface.IMsgHandler
	msgChan        chan []byte
	ExitBufferChan chan bool
}

func NewConnection(conn *net.TCPConn, connId uint32, MsgHandler ziface.IMsgHandler) *Connection {
	return &Connection{
		Conn:           conn,
		ConnId:         connId,
		IsClose:        false,
		msgChan:        make(chan []byte),
		ExitBufferChan: make(chan bool),
		MsgHandler:     MsgHandler,
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()
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
		req := &Request{
			conn: c,
			data: msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

func (c *Connection) StartWriter() {
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.GetTcpConnection().Write(data); err != nil {
				c.ExitBufferChan <- true
			}
		case <-c.ExitBufferChan:
			fmt.Println("Reader Goroutine exit")
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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClose == true {
		return errors.New("client has already closed")
	}
	dp := NewDataPack()
	dataPack, err := dp.Pack(NewMessage(data, msgId))
	if err != nil {
		return err
	}
	c.msgChan <- dataPack
	return nil
}
