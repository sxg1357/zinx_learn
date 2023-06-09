package znet

import "zinx_learn/ziface"

type Request struct {
	conn ziface.IConnection
	data []byte
}

func (r *Request) GetData() []byte {
	return r.data
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}
