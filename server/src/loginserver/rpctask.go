package main

import (
	"base/env"
	"github.com/golang/protobuf/proto"
	"net/rpc"
	"common"
	"encoding/json"
	"fmt"
	"glog"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"usercmd"
)

type RPCLogicTask int

func (this *RPCLogicTask) VerifyHandle()bool {

	return false
}

func (this *RPCLogicTask) OnClose(){

}


func (this *RPCLogicTask) OnRecv(conn , id uint64, module uint8, cmd uint16, data []byte)(/*grpc.PbObj*/[]byte,  int) {

	data = append(data, byte(id), byte(id>>8), byte(id>>16), byte(id>>24), byte(id>>32), byte(id>>40), byte(id>>48), byte(id>>56))
	header := make(http.Header)
	header.Set("qqt","1")
	req := &http.Request{
		Header: header,
		ContentLength: int64(len(data)),
		Body: &common.HttpBody{
			Buf: data,
		},
	}
	res := &common.ResWrite{}

	switch usercmd.SRPCLogin(cmd) {
	case usercmd.SRPCLogin_Logingame:
		HandleGameMsg(res, req)
	case usercmd.SRPCLogin_Loginlogin:
		HandleLoginMsg(res, req)
	//case usercmd.SRPCLogin_LoginServer:
	//	HandleServer(id, data)
	}
	l := len(res.Buf)
	if l == 0{
		return  nil,  1
	}
	return common.RetMsg(res.Buf),  0
}

///////
func (handler *RPCLogicTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	raddr := env.Get("rcenter", "listen")
	client, err := rpc.Dial("tcp", raddr)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer func() {
		if err := client.Close();err == nil {
			glog.Error("[RPC] Client close error", err)
		}
	}()

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		glog.Error("[RPC] Read request.Body failed, err:%v\n", err)
		return
	}

	data:= make(map[string]string)
	err = json.Unmarshal(body, &data)
	if err !=nil {
		glog.Error("[RPC] Body To Json failed, err:", err)
		panic("Error")
	}
	fmt.Println("++++++++++++++",data)

	//w.Write([]byte(data["username"]))


	uid, _ := strconv.Atoi(data["uid"])
	var reply = &usercmd.RetIntoFRoom{}
	var param = &usercmd.ReqIntoRoom{
		UId: proto.Uint64(uint64(uid)),
		UserName: proto.String(data["username"]),
	}
	if reply.GetKey() == ""{
		fmt.Println("------------",uid, "--" )
	}


	err = client.Call("RetIntoRoom.RetRoom", &param, &reply)
	fmt.Println("-----Call end-------",uid, "--" )

	if err != nil {
		log.Fatal(err)
	}
	if reply.GetKey() == ""{
		fmt.Println("----////--",uid, "--" )
	}

	w.Write([]byte(reply.GetKey()))
}







////*///// gRPC server
type RoomService struct {

}

func (this *RoomService) Route(conn usercmd.StreamRoomService_RouteServer) error {
	for {
		stream, err := conn.Recv()
		if io.EOF == err {
			glog.Info("[gRPC] Server Got EOF")
			return nil
		}
		if err != nil {
			glog.Error("[gRPC] Server Error. ", err)
			return err
		}
		glog.Info("[gRPC] Server Recv: ", stream.Data)
		switch stream.Type {
		case usercmd.MsgTypeRpc_Regist:
			var info usercmd.ConnectRoomInfo
			err := json.Unmarshal(stream.Data, &info)
			if err != nil {
				glog.Error("[Common] Json to struct error.", err)
				return err
			}
			fmt.Println("Server got regist msg ", info.Ip, ", ", &info.Port)
			// Add room
			// to do
			break
		case usercmd.MsgTypeRpc_Update:
			break

		}

	}
}

func StartGRPCServer() bool {

	addr := env.Get("rcenter", "grpc")
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Error("[Start] Binding port error, ", addr, ", ", err)
		return false
	}
	s := grpc.NewServer()
	usercmd.RegisterStreamRoomServiceServer(s, &RoomService{})
	go func() {
		s.Serve(listen)
	}()
	glog.Info("[gRPC] Start server success, ", s.GetServiceInfo())
	return true

}