package main

import (
	"glog"
	"net"
	"net/http"
	"time"

)

func startHttpClient(){
	x := http.NewServeMux()
	x.HandleFunc("/getroom", GetRoom)
	//x.HandleFunc("/getload", GetLoad)
	listen, err := net.Listen("tcp",":9090")
	if err != nil {
		glog.Error("Binding Failed.")
	}
	ser := &http.Server{
		WriteTimeout:  60 * time.Second,
		ReadTimeout:   60 * time.Second,
		Handler: x,
	}
	go ser.Serve(listen)
	glog.Info("Http binding successful.")

}

func GetRoom(w http.ResponseWriter, req *http.Request){
	//id, _ := strconv.ParseUint(req.URL.Query()["id"][0],10,32)


}