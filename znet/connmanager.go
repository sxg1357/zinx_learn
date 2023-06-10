package znet

import (
	"fmt"
	"sync"
	"zinx_learn/ziface"
)

type ConnManager struct {
	connection map[uint32]ziface.IConnection
	mutex      sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connection: make(map[uint32]ziface.IConnection),
	}
}

func (cn *ConnManager) Add(conn ziface.IConnection) {
	defer cn.mutex.Unlock()
	cn.mutex.Lock()
	cn.connection[conn.GetConnId()] = conn
	fmt.Println("connection add to connManager successfully, connId:", conn.GetConnId(), " len:", cn.Len())
}

func (cn *ConnManager) Remove(conn ziface.IConnection) {
	defer cn.mutex.Unlock()
	cn.mutex.Lock()
	if _, ok := cn.connection[conn.GetConnId()]; !ok {
		fmt.Println("connId:", conn.GetConnId(), " not exists in connectionManager")
		return
	}
	delete(cn.connection, conn.GetConnId())
}

func (cn *ConnManager) Len() uint32 {
	return uint32(len(cn.connection))
}

func (cn *ConnManager) ClearConn() {
	defer cn.mutex.Unlock()
	cn.mutex.Lock()
	for connId, conn := range cn.connection {
		//先将连接关闭
		conn.Stop()
		delete(cn.connection, connId)
	}
}
