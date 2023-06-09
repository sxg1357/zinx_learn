package znet

import (
	"errors"
	"fmt"
	"net"
	"time"
	"zinx_learn/ziface"
)

type Server struct {
	ServerName string
	IpVersion  string
	Ip         string
	Port       int
}

func NewServer(name string) ziface.IServer {
	return &Server{
		ServerName: name,
		IpVersion:  "tcp4",
		Ip:         "0.0.0.0",
		Port:       9501,
	}
}

func CallBakClientApi(conn *net.TCPConn, data []byte, cnt int) error {
	//将客户端的数据回显
	fmt.Println("[Conn Handle] CallBackToClient ... ")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallBakClientApi error")
	}
	return nil
}

func (s *Server) Start() {
	fmt.Printf("[%s] start listening on %s:%d\r\n", s.ServerName, s.Ip, s.Port)

	addr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("ResolveTCPAddr failed...")
		return
	}

	go func() {
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
				dealConn := NewConnection(conn, cid, CallBakClientApi)
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
