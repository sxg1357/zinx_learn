package ziface

type IConnManager interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	Len() uint32
	ClearConn()
}
