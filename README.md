# zinx_learn

Server.go
```go
package main

import (
	"fmt"
	"zinx_learn/ziface"
	"zinx_learn/znet"
)

type MyRoute struct {
	znet.BaseRouter
}

func (mr *MyRoute) PreHandler(request ziface.IRequest) {
	fmt.Println("PreHandler execute...")
	fmt.Printf("recv %s from client\r\n", string(request.GetData()))
}

func (mr *MyRoute) Handler(request ziface.IRequest) {
	fmt.Println("Handler execute...")
	//if err := c.HandleFunc(c.Conn, buf, cnt); err != nil {
	//	fmt.Println("call Handler Api error...")
	//	c.ExitBufferChan <- true
	//	return
	//}
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("write error")
	}
}

func (mr *MyRoute) PostHandler(request ziface.IRequest) {
	fmt.Println("PostHandler execute...")
}

func main() {
	s := znet.NewServer()
	s.AddRouter(&MyRoute{})
	s.Server()
}

```

Client.go

```go
package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("client test start...")
	conn, err := net.Dial("tcp4", "127.0.0.1:9501")
	if err != nil {
		fmt.Println("connect server error")
	}

	for {
		conn.Write([]byte("Hello, World"))
		buf := make([]byte, 512)
		cnts, err := conn.Read(buf)
		if cnts == 0 {
			fmt.Println("server close")
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("read error")
			continue
		}
		fmt.Printf("recv %d bytes from server, msg:%s\r\n", cnts, string(buf[:cnts]))
		time.Sleep(time.Second * 2)
	}
}
```

zinx.json

```json
{
  "Name":"zinx",
  "Host":"0.0.0.0",
  "Port":9501,
  "MaxConn":9999
}
```

