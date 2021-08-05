package main

import (
	"base/env"
	"base/gonet"
	"flag"
	"glog"
	"math/rand"
	"net/http"
	"time"
)

const (
	TokenRedis int = iota
)

type RoomServer struct{
	gonet.Service
	roomser 	*gonet.TcpServer
	//roomserUdp 	*snet.Server
	version uint32
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
	glog.Info("[Start] Initialization.")

	//check
	pprofport := env.Get("room", "port")
	if pprofport != "" {
		go func() {
			http.ListenAndServe(pprofport, nil)
		}()
	}

	//Global config
	//if()

	//Redis
	//To do


	// Binding Local Port
	err := this.roomser.Bind(env.Get("room", "listen"))
	if err != nil {
		glog.Error("[Start] Binding port failed")
		return false
	}


	//
	if !RCenterClient_GetMe().Connect(){
		return false
	}

	glog.Info("[Start] Initialization successful, ", this.version)
	return true
}

func (this *RoomServer) UdpLoop() {
	//for {
	//	//conn, err := this.roomserUdp.Accept()
	//}
}

func (this *RoomServer) MainLoop() {
	conn, err := this.roomser.Accept()
	if err != nil {
		return
	}
	NewPlayerTask(conn).Start()
}

func (this *RoomServer) Final() bool {
	this.roomser.Close()


	return true
}

func (this *RoomServer) Reload() {

}

var (
	logfile = flag.String("logfile", "","Log file name")
	config = flag.String("config", "config.json","config path")
)

func main() {
	flag.Parse()

	if !env.Load(*config){
		return
	}

	//loglevel := env.Get("global", "loglevel")
	//if loglevel != "" {
	//	flag.Lookup("stderrthreshold").Value.Set(loglevel)
	//}
	//
	//logtostderr := env.Get("global", "logtostderr")
	//if loglevel != "" {
	//	flag.Lookup("logtostderr").Value.Set(logtostderr)
	//}

	rand.Seed(time.Now().Unix())

	if *logfile != ""{
		glog.SetLogFile(*logfile)
	}else{
		glog.SetLogFile(env.Get("room","log"))
	}

	defer glog.Flush()

	RoomServer_GetMe().Main()

	glog.Info("[Close] RoomServer closed.")
}
