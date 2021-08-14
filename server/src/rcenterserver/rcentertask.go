package main

import (
	"base/gonet"
	"net"
	"sync"
	"usercmd"
)

type ServerTask struct {
	gonet.TcpTask
	Id      uint16
	Addr    string
	NewSync bool
	RoomNum uint32
	UserNum uint32
	weight  uint32 // Load weight?
	index   int
}

func NewServerTask(conn net.Conn) *ServerTask {
	s := &ServerTask{
		TcpTask: *gonet.NewTcpTask(conn),
	}
	s.Derived = s
	return s
}

func (this *ServerTask) Index() int {
	return this.index
}

func (this *ServerTask) SetIndex(i int) {
	this.index = i
}

func (this *ServerTask) Weight() uint32 {
	return this.weight
}

func (this *ServerTask) SetWeight(w uint32) {
	this.weight = w
}

func (this *ServerTask) Key() string {
	return this.Addr
}

func (this *ServerTask) ParseMsg(data []byte, flag byte) bool {

	return true
}

func (this *ServerTask) OnClose() {

}

func (this *ServerTask) SendCmd(cmd usercmd.CmdType, msg []byte) bool {

	return true
}

// Server task manager
type ServerTaskMgr struct {
	mutex    sync.RWMutex
	servers  map[uint16]*ServerTask
	uniqueid uint32
	//serlist 	*ServerList
	tnumtime int64
	totalnum uint32
}

var servertm *ServerTaskMgr

func ServerTaskMgr_GetMe() *ServerTaskMgr {
	if servertm == nil {
		servertm = &ServerTaskMgr{
			servers: make(map[uint16]*ServerTask),
		}
	}
	return servertm
}

func (this *ServerTaskMgr) GetServer() (uint16, string, bool) {

	// -----------------xie si--------
	// To do
	// -----------------xie si--------
	return uint16(10000), "192.168.96.72", true
}
