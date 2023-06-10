package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx_learn/utils"
	"zinx_learn/ziface"
)

type Connection struct {
	Conn           *net.TCPConn
	ConnId         uint32
	IsClose        bool
	MsgHandler     ziface.IMsgHandler
	msgChan        chan []byte
	msgChanBuffer  chan []byte
	ExitBufferChan chan bool
	TcpServer      ziface.IServer
	Property       map[string]interface{}
	PropertyLock   sync.RWMutex
}

func NewConnection(conn *net.TCPConn, connId uint32, MsgHandler ziface.IMsgHandler, server ziface.IServer) *Connection {
	c := &Connection{
		Conn:           conn,
		ConnId:         connId,
		IsClose:        false,
		msgChan:        make(chan []byte),
		msgChanBuffer:  make(chan []byte, utils.GlobalObject.MaxPacketLen),
		ExitBufferChan: make(chan bool),
		MsgHandler:     MsgHandler,
		TcpServer:      server,
		Property:       make(map[string]interface{}),
	}
	conMgr := c.TcpServer.GetConnMgr()
	conMgr.Add(c)
	c.TcpServer.CallConnOnStart(c)
	c.SetProperty("name", "sxg")
	c.SetProperty("age", 25)
	return c
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
		case data, ok := <-c.msgChanBuffer:
			if ok {
				if _, err := c.GetTcpConnection().Write(data); err != nil {
					c.ExitBufferChan <- true
				}
			} else {
				fmt.Println("msgBufferChan is closed")
				break
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
	c.TcpServer.GetConnMgr().Remove(c)
	c.IsClose = true
	c.ExitBufferChan <- true
	c.Conn.Close()
	c.TcpServer.CallConnOnStop(c)
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

func (c *Connection) SetProperty(key string, val interface{}) {
	defer c.PropertyLock.Unlock()
	c.PropertyLock.Lock()
	c.Property[key] = val
}

func (c *Connection) GetProperty(key string) interface{} {
	defer c.PropertyLock.RUnlock()
	c.PropertyLock.RLock()
	if val, ok := c.Property[key]; ok {
		return val
	} else {
		fmt.Println("property not found")
		return nil
	}
}

func (c *Connection) RemoveProperty(key string) {
	defer c.PropertyLock.Unlock()
	c.PropertyLock.Lock()
	if _, ok := c.Property[key]; ok {
		delete(c.Property, key)
		fmt.Println("remove property ", key, " ok")
	} else {
		fmt.Println("property not found")
	}
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
	c.msgChanBuffer <- dataPack
	return nil
}
