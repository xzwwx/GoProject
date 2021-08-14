package main

import (
	"base/env"
	"base/gonet"
	"common"
	"flag"
	"fmt"
	"glog"
	"math/rand"
	"net"
	"time"
)

const (
	TokenRedis int = iota
)

type RoomServer struct {
	gonet.Service
	roomser *gonet.TcpServer
	//roomserUdp 	*snet.Server
	version uint32
}

var serverm *RoomServer

func RoomServer_GetMe() *RoomServer {
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
	fmt.Println(pprofport)
	if pprofport != "" {
		go func() {
			fmt.Println("开始监听：", pprofport)
			net.Listen("tcp", pprofport)
			//http.ListenAndServe(pprofport, nil)
		}()
	}

	if !common.RedisMgr.NewRedisManager() {
		glog.Errorln("[LogicServer Init] Init error")
		return false
	}
	//Global config
	//if()

	//Redis
	//To do

	fmt.Println("-------：", pprofport)

	// Binding Local Port
	err := this.roomser.Bind(env.Get("room", "listen"))
	if err != nil {
		glog.Error("[Start] Binding port failed")
		return false
	}

	//
	//if !RCenterClient_GetMe().Connect(){
	//return false
	//}

	glog.Info("[Start] Initialization successful, ", this.version)
	return true
}

func (this *RoomServer) UdpLoop() {
	//for {
	//	//conn, err := this.roomserUdp.Accept()
	//}
}

func ClientLogic(conn net.Conn) {

	// 从客户端接受数据
	buf := make([]byte, 1024)

	n, _ := conn.Read(buf)
	//s, _ := bufio.NewReader(conn).Read(buf)
	println("由客户端发来的消息：", n, ", ", string(buf[:n]))

	// 发送消息给客户端
	//conn.Write([]byte("东东你好\n"))

	// 关闭连接
	//conn.Close()
}

func (this *RoomServer) MainLoop() {
	// fmt.Println("loop")

	conn, err := this.roomser.Accept()
	if err != nil {
		return
	}
	// go ClientLogic(conn)
	NewPlayerTask(conn).Start()
}

func (this *RoomServer) Final() bool {
	this.roomser.Close()

	return true
}

func (this *RoomServer) Reload() {

}

var (
	logfile = flag.String("logfile", "", "Log file name")
	config  = flag.String("config", "config.json", "config path")
)

func main() {
	flag.Parse()

	if !env.Load(*config) {
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

	if *logfile != "" {
		glog.SetLogFile(*logfile)
	} else {
		glog.SetLogFile(env.Get("room", "log"))
	}

	defer glog.Flush()
	fmt.Println("logok")
	RoomServer_GetMe().Main()

	glog.Info("[Close] RoomServer closed.")
}
