package main

import (
	"common"
	"net/http"
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