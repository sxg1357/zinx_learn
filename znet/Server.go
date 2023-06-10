package znet

import (
	"fmt"
	"net"
	"time"
	"zinx_learn/utils"
	"zinx_learn/ziface"
)

type Server struct {
	ServerName string
	IpVersion  string
	Ip         string
	Port       int
	MsgHandler ziface.IMsgHandler
	ConnMgr    ziface.IConnManager
}

func NewServer() ziface.IServer {
	//utils.GlobalObject.Reload()
	return &Server{
		ServerName: utils.GlobalObject.Name,
		IpVersion:  "tcp4",
		Ip:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.Port,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
}

//func CallBakClientApi(conn *net.TCPConn, data []byte, cnt int) error {
//	//将客户端的数据回显
//	fmt.Println("[Conn Handle] CallBackToClient ... ")
//	if _, err := conn.Write(data[:cnt]); err != nil {
//		fmt.Println("write back buf err ", err)
//		return errors.New("CallBakClientApi error")
//	}
//	return nil
//}

func (s *Server) Start() {
	fmt.Printf("server start listening on %s:%d\r\n", s.Ip, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketLen: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketLen)

	addr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("ResolveTCPAddr failed...")
		return
	}

	go func() {
		go s.MsgHandler.StartWorkerPool()
		listener, err := net.ListenTCP(s.IpVersion, addr)
		if err != nil {
			fmt.Println("ListenTCP failed...")
			return
		}

		var cid uint32
		cid = 0
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP failed...")
				continue
			}
			//接收客户端连接开启一个协程处理消息
			go func() {
				dealConn := NewConnection(conn, cid, s.MsgHandler, s)
				cid++

				//启动一个协程处理当前连接的业务
				go dealConn.Start()
			}()
		}
	}()
}

func (s *Server) Server() {
	s.Start()
	for {
		time.Sleep(time.Second * 10)
	}
}

func (s *Server) Stop() {
	fmt.Printf("[Stop] %s server...\r\n", s.ServerName)
}

func (s *Server) AddRouter(msgId uint32, route ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, route)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}
