package main

import (
	"base/env"
	//"base/env"
	"base/rpc"
	"time"

	"base/gonet"
	//"bytes"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"glog"
	"strconv"

	//"google.golang.org/genproto/googleapis/ads/googleads/v3/common"
	"log"
	"net"
	"usercmd"
)

type RCenterServer struct {
	gonet.Service
	rpcser 	*gonet.TcpServer
	sockser	*gonet.TcpServer
}

var serverm *RCenterServer

func RCenterServer_GetMe() *RCenterServer {
	if serverm == nil {
		serverm = &RCenterServer{
			rpcser:  &gonet.TcpServer{},
			sockser:  &gonet.TcpServer{},
		}
		serverm.Derived = serverm
	}
	return serverm
}

///////rpc
type RetIntoRoom struct {
}

func (q *RetIntoRoom) RetRoom(request *usercmd.ReqIntoRoom, reply *usercmd.RetIntoFRoom) error{

	fmt.Println("Into RPC...")

	uid := request.GetUId()
	username := request.UserName
	fmt.Println(strconv.FormatInt(int64(uid), 10))


	// Request into room
	// To do



	reply.Err = proto.Uint32(uint32(0))
	reply.RoomId = proto.Uint32(uint32(107))
	reply.Addr = proto.String("127.0.0.1:9494")
	reply.Key = proto.String(strconv.FormatInt(int64(uid), 10)+*username)

	fmt.Println("Out RPC")

	return nil
}

func Acc(listener net.Listener){
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		fmt.Println("Accept success...")
		go rpc.ServeConn(conn)
		fmt.Println("Server success...")
	}
}
/////////

func (this * RCenterServer) Init() bool{


	rpc.RegisterName("RPCTask", new(RPCTask))
	fmt.Println("RegisterName success...")

	addr := env.Get("rcenter", "listen")
	listener, err := net.Listen("tcp", addr)
	fmt.Println("Listen success...")

	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	go Acc(listener)

	fmt.Println("Init success.")
	return true
}

func (this *RCenterServer) MainLoop() {
	time.Sleep(time.Second)
}

func (this *RCenterServer) Final() bool {

	return true
}

func (this *RCenterServer) Reload() {

}

var (
	logfile = flag.String("logfile", "","Log file name")
	config = flag.String("config", "config.json","config path")
)

func main() {
	flag.Parse()
	//
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
	//
	if *logfile != ""{
		glog.SetLogFile(*logfile)
	}else{
		glog.SetLogFile(env.Get("rcenter","log"))
	}

	defer glog.Flush()

	RCenterServer_GetMe().Main()

	glog.Info("[Close] RCenterServer closed.")
}

