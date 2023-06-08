package znet

import (
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

func (s *Server) Start() {
	fmt.Printf("[%s v0.1] start listening on %s:%d", s.ServerName, s.Ip, s.Port)

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

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP failed...")
				continue
			}
			//接收端哦客户端连接开启一个协程处理消息
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv msg err...")
						continue
					}
					fmt.Printf("recv %d bytes from client, msg:%s", cnt, string(buf))
				}
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
	fmt.Printf("[Stop] %s server...", s.ServerName)
}

func NewServer(name string) ziface.IServer {
	return &Server{
		ServerName: name,
		IpVersion:  "tcp4",
		Ip:         "0.0.0.0",
		Port:       9501,
	}
}
