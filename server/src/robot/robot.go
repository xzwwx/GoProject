package main

import (
	"common"
	glog "glog-master"
	"gonet"
	"net/url"
	"strconv"
	"strings"
	"usercmd"
)

const (
	rpc_header_size = 4
	req_header_size = 11
	ret_header_size = 8
)

type LogicClient struct {
	RpcTask
	mclient 	*gonet.TcpClient
	addr 		string
	key 		string
	shutdown 	bool
	//aesdec
	//aesenc
	//privKey
	tmpseqid	uint32
	connected 	bool
	id uint64

	loginInfo 	LoginInfo
	//teamcli		*TeamClient
	//roomcli		*RoomClient

}

func NewLogicClient(addr, key string, id uint64, name string) *LogicClient{
	sclient := &LogicClient{
		RpcTask:	*NewRpcTask(nil),
		mclient: 	&gonet.TcpClient{},
		addr: 		addr,
		key:		key,
		//privKey:
		//
		tmpseqid: 	1,
		connected: 	false,
		id:			id,
		//teamcli:	NewTeamClient(),
		//	roomcli:	NewRoomClient(),
	}
	sclient.Drived = sclient
	//
	//
	return sclient
}

func (this * LogicClient)SendCmd(mainCmd uint8, subCmd uint16, SeqId uint32, Buff []byte) bool {

	// simplified
	rSize :=len(Buff)
	tmpinbuff := make([]byte, rSize)
	tmpinbuff[4] = byte(mainCmd)
	tmpinbuff[5] = byte(subCmd>>8)
	tmpinbuff[6] = byte(subCmd)


	msglen := rpc_header_size +len(tmpinbuff)
	sendBuff := make([]byte, msglen)

	copy(sendBuff[rpc_header_size:], tmpinbuff)

	return this.SendBytes(sendBuff)

}


//send RPC message
func (this *LogicClient) SendRpcCmd(uri string)bool{
	args := strings.Split("rui", "?")
	if len(args) != 2{
		glog.Error("[Logic] parameters error.", uri)
		return false
	}

	var (
		mainCmd uint8
		subCmd	uint16
	)
	switch args[0]{
	case "/msg":
		mainCmd = common.MsgType_Login
		subCmd = uint16(usercmd.SRPCLogin_Loginmsg)
	case "/game":
		mainCmd = common.MsgType_Login
		subCmd = uint16(usercmd.SRPCLogin_Logingame)
	default:
		glog.Error("[Logic] Protocal error.")
		return false
	}

	reqCmd := &usercmd.ReqHttpArgData{}
	values, err := url.ParseQuery(args[1])
	if err != nil{
		return false
	}
	for key, value := range values{
		if key == "c"{
			reqCmd, _ = strconv.Atoi(value[0])
			reqCmd.Cmd = uint32(reqCmd)
		} else {
			args := &usercmd.ReqHttpArgData_KeyVal{
				Key: key,
				Val: value[0],
			}
			reqCmd.Args = append(reqCmd.Args, args)
		}
	}
	data, err := reqCmd.Marshal()
	if err != nil{
		return false
	}
	this.tmpseqid ++
	return this.SendCmd(mainCmd, subCmd, this.tmpseqid, data)
}

