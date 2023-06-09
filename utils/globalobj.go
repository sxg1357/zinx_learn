package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx_learn/ziface"
)

type GlobalObj struct {
	TcpServer    ziface.IServer
	Name         string
	Host         string
	Port         int
	Version      string
	MaxConn      uint32
	MaxPacketLen uint32
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:         "zinx",
		Host:         "0.0.0.0",
		Port:         9501,
		Version:      "v0.3",
		MaxConn:      10240,
		MaxPacketLen: 4096,
	}
	GlobalObject.Reload()
}

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("config/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
