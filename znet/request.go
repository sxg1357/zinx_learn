package znet

import "zinx_learn/ziface"

type Request struct {
	conn ziface.IConnection
	//data []byte
	data ziface.IMessage
}

func (r *Request) GetData() []byte {
	return r.data.GetData()
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}
