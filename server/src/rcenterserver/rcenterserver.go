package main

import (
	"base/env"
	"flag"
	"gonet"
	"glog"
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

func (this * RCenterServer) Init() bool{


	return true
}

func (this *RCenterServer) MainLoop() {

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

	if !env.Load(*config){
		return
	}
	loglevel := env.Get("global", "loglevel")
	if loglevel != "" {
		flag.Lookup("stderrthreshold").Value.Set(loglevel)
	}

	logtostderr := env.Get("global", "logtostderr")
	if loglevel != "" {
		flag.Lookup("logtostderr").Value.Set(logtostderr)
	}

	if *logfile != ""{
		glog.SetLogFile(*logfile)
	}else{
		glog.SetLogFile(env.Get("rcenter","log"))
	}

	defer glog.Flush()

	RCenterServer_GetMe().Main()

	glog.Info("[Close] RCenterServer closed.")
}