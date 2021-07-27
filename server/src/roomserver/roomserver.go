package main

import "gonet"

const (
	TokenRedis int = iota
)

type RoomServer struct{
	gonet.Service
	roomser 	*gonet.TcpServer
	//roomserUdp 	*snet.Server
}

var serverm *RoomServer

func RoomServer_GetMe() *RoomServer{
	if serverm == nil {
		serverm = &RoomServer{
			roomser: &gonet.TcpServer{},
			//roomerUdp
		}
		serverm.Derived = serverm
	}
	return serverm
}

func (this *RoomServer) Init() bool {



	return true
}

func (this *RoomServer) UdpLoop() {

}

func (this *RoomServer) MainLoop() {

}

func (this *RoomServer) Final() bool {


	return true
}

func (this *RoomServer) Reload() {

}

