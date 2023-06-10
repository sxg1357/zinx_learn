# zinx_learn

Server.go
```go
package main

import (
	"fmt"
	"zinx_learn/ziface"
	"zinx_learn/znet"
)

type MyRoute1 struct {
	znet.BaseRouter
}

type MyRoute2 struct {
	znet.BaseRouter
}

func (mr *MyRoute1) PreHandler(request ziface.IRequest) {
	fmt.Println("PreHandler1 execute...")
	fmt.Println("MsgId:", request.GetMsgId(), " data:", string(request.GetData()))
}

func (mr *MyRoute1) Handler(request ziface.IRequest) {
	fmt.Println("Handler1 execute...")
	if err := request.GetConnection().SendMsg(request.GetMsgId(), []byte("ping...ping...ping")); err != nil {
		fmt.Println(err)
	}
}

func (mr *MyRoute2) PostHandler(request ziface.IRequest) {
	fmt.Println("PostHandler1 execute...")
}

func (mr *MyRoute2) PreHandler(request ziface.IRequest) {
	fmt.Println("PreHandler2 execute...")
	fmt.Println("MsgId:", request.GetMsgId(), " data:", string(request.GetData()))
}

func (mr *MyRoute2) Handler(request ziface.IRequest) {
	fmt.Println("Handler2 execute...")
	if err := request.GetConnection().SendMsg(request.GetMsgId(), []byte("ping...ping...ping")); err != nil {
		fmt.Println(err)
	}
}

func (mr *MyRoute1) PostHandler(request ziface.IRequest) {
	fmt.Println("PostHandler2 execute...")
}

func main() {
	s := znet.NewServer()
	s.AddRouter(1, &MyRoute1{})
	s.AddRouter(2, &MyRoute2{})
	s.Server()
}
```

Client.go

```go
package main

import (
	"fmt"
	"net"
	"zinx_learn/znet"
)

func main() {
	fmt.Println("client test start...")
	conn, err := net.Dial("tcp4", "127.0.0.1:9501")
	if err != nil {
		fmt.Println("connect server error")
	}

	//for {
	//	conn.Write([]byte("Hello, World"))
	//	buf := make([]byte, 512)
	//	cnts, err := conn.Read(buf)
	//	if cnts == 0 {
	//		fmt.Println("server close")
	//		return
	//	}
	//	if err != nil && err != io.EOF {
	//		fmt.Println("read error")
	//		continue
	//	}
	//	fmt.Printf("recv %d bytes from server, msg:%s\r\n", cnts, string(buf[:cnts]))
	//	time.Sleep(time.Second * 2)
	//}

	dp := znet.NewDataPack()
	msg1 := &znet.Message{
		Data:   []byte{'h', 'e', 'l', 'l', 'o'},
		MsgLen: 5,
		MsgId:  1,
	}

	packData1, err := dp.Pack(msg1)
	if err != nil {
		return
	}

	_, err = conn.Write(packData1)
	if err != nil {
		return
	}

	msg2 := &znet.Message{
		Data:   []byte{'w', 'o', 'r', 'l', 'd'},
		MsgLen: 5,
		MsgId:  1,
	}

	packData2, err := dp.Pack(msg2)
	if err != nil {
		return
	}

	_, err = conn.Write(packData2)
	if err != nil {
		return
	}

	dataPack3 := append(packData1, packData2...)
	_, err = conn.Write(dataPack3)
	if err != nil {
		return
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


ServerPack.go

```go
package main

import (
	"fmt"
	"io"
	"net"
	"zinx_learn/znet"
)

func main() {
	listener, err := net.Listen("tcp4", "0.0.0.0:9501")
	if err != nil {
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go func(conn net.Conn) {
			for {
				dp := znet.NewDataPack()
				headData := make([]byte, dp.GetHeadLen())

				_, err := io.ReadFull(conn, headData)
				if err != nil {
					return
				}

				msgHead, err := dp.UnPack(headData)
				if err != nil {
					return
				}

				if msgHead.GetDataLen() > 0 {
					msg := msgHead.(*znet.Message)
					msg.Data = make([]byte, msg.GetDataLen())
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						return
					}
					fmt.Printf("Recv Msg:[%s], MsgId:[%d], MsgLen:%d\r\n", msg.Data, msg.MsgId, msg.MsgLen)
				}
			}
		}(conn)
	}

}
```



DataPackClient.go
```go
package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx_learn/znet"
)

func main() {
	fmt.Println("client test start...")
	conn, err := net.Dial("tcp4", "127.0.0.1:9501")
	if err != nil {
		fmt.Println("connect server error")
	}

	for {
		dp := znet.NewDataPack()
		//发给路由1
		msg1 := znet.NewMessage([]byte("Hello"), 1)
		packData1, _ := dp.Pack(msg1)
		conn.Write(packData1)

		//发给路由2
		msg2 := znet.NewMessage([]byte("World"), 2)
		packData2, _ := dp.Pack(msg2)
		conn.Write(packData2)

		//解析数据包
		headData := make([]byte, dp.GetHeadLen())
		cnts, err := io.ReadFull(conn, headData)
		if cnts == 0 {
			fmt.Println("server close")
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("read error")
			continue
		}
		unpackData, err := dp.UnPack(headData)
		if err != nil {
			return
		}
		if unpackData.GetDataLen() > 0 {
			data := unpackData.(*znet.Message)
			data.Data = make([]byte, data.GetDataLen())
			_, err = io.ReadFull(conn, data.Data)
			if err != nil {
				return
			}
			fmt.Println("msg:", string(data.Data), " msgId:", data.MsgId)
		}
		//fmt.Printf("recv %d bytes from server, msg:%s\r\n", cnts, string(buf[:cnts]))
		time.Sleep(time.Second * 2)
	}
}

```




