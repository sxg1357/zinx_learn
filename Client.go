package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("test client start...")
	time.Sleep(time.Second * 3)

	conn, err := net.Dial("tcp4", "127.0.0.1:9501")
	if err != nil {
		fmt.Println("connect server error")
		return
	}

	for {
		_, err := conn.Write([]byte("Hello, World"))
		if err != nil {
			fmt.Println("write error")
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error...")
			return
		}
		fmt.Printf("recv %d bytes from server, msg:%s", cnt, string(buf))
		time.Sleep(time.Second * 1)
	}

}
