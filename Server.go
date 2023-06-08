package main

import "zinx_learn/znet"

func main() {
	s := znet.NewServer("zinx")
	s.Server()
}
