package main

import (
	"common"
	glog "glog-master"
	"gonet"
	"net/url"
	"strconv"
	"strings"
	"time"
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
	roomcli		*RoomClient

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

func (this *LogicClient) Connect() bool {
	conn, err := this.mclient.Connect(this.addr)
	if err != nil {
		glog.Error("[Logic] Connect failed.", this.addr)
		return false
	}
	this.Conn = conn
	this.Start()
	this.Verify()


	//////////////////// Connect Server could be modified
	reqCmd := &usercmd.ReqGateLogin{
		Key: 	this.key,
		Version:"1.0",
		StrArgs:this.privKey.N.Text(16),
		IntArgs:uint32(this.privKey.E),
	}

	reqData := make([]byte, reqCmd.Size())
	_, err = reqCmd.MarshalTo(reqData)
	if err != nil {
		glog.Error("[Protocal] Encode failed.")
		return false
	}
	//////////////////////////
	this.tmpseqid++
	this.SendCmd(common.MsgType_LogicGate, uint16(usercmd.GateCmd_GateLogin), this.tmpseqid, reqData)

	glog.Infof("[Logic] Connect server success.", this.addr, ",", this.key)
	return true
}


func (this *LogicClient) ParseMsg(data []byte) bool{
	if len(data) < 4{
		glog.Error("[Logic] Message error.", this.Conn.RemoteAddr(), ",", data)
		return false
	}

	// decrypt
	//if (data[3] & )

	return true
}

func (this *LogicClient) OnClose() {
	this.Reset()

	glog.Info("[Logic] Disconnect with ", this.addr, ", ", this.shutdown)
	for !this.shutdown{
		glog.Info("[Logic] Reconnecting..", this.addr)
		if this.Connect(){
			glog.Info("[Logic] Server reconnected.", this.addr)
			break
		}
		time.Sleep(time.Second * 3)
	}
}






////// Simplified mode
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



func (this *LogicClient) ParseRpcMsg(data []byte) bool {
	httpflag := data[3]
	httpcmd := usercmd.MsgType(int(data[5])<<8 | int(data[4]))

	glog.Info("[Logic] Receive message: ", httpflag, "cmd=", httpcmd, ", len=",len(data))

	switch httpcmd {
	case usercmd.MsgType_ReqIntoFRoom:			/////////////
		revCmd, ok := common.DecodeCmd(data[4:], httpflag, &usercmd.RetIntoFRoom{}).(*usercmd.RetIntoFRoom)				///////
		if !ok {
			return false
		}

		if *revCmd.Err == common.ErrorCodeUseFreeMatch {
			glog.Info("[Login] Team match.")
			return false
		}

		if *revCmd.Err == common.ErrorCodeOkay {
			glog.Error("[Login] Room request error.", *revCmd.Err)
			return false
		}

		this.loginInfo.Address = *revCmd.Addr
		this.loginInfo.Key = *revCmd.Key
		glog.Info("[Login] Room success.", this.loginInfo.Address)

		this.roomcli.Connect(this.loginInfo.Address, this.loginInfo.Key, "test")
		glog.Info("[Login] Room connected.",this.loginInfo.Address)

	}
}

// Room Client
type RoomClient struct {
	gonet.TcpTask
	roomcli 	*gonet.TcpClient
	id 			uint64
	name		string
	dev			string
}

func NewRoomClient(name string) *RoomClient{
	roomcli := &RoomClient{
		TcpTask:	*gonet.NewTcpTask(nil),
		name:		name,
	}
	roomcli.Derived = roomcli
	return roomcli
}

func (this *RoomClient) Connect(address, key, name string) bool {
	conn, err := this.roomcli.Connect(address)
	if err != nil{
		glog.Error("[Player] Connection failed.", address)
		return false
	}
	this.Conn = conn
	this.Start()
	this.SendCmd(usercmd.MsgTypeCmd_Login, &usercmd.MsgLogin{
		Name : name,
		key : key,
	})
	return true
}

func (this *RoomClient) ParseMsg(data []byte, flag byte) bool{


	return true
}

func (this *RoomClient) OnClose() {

}

func (this *RoomClient) SendCmd( cmd usercmd.MsgTypeCmd, msg usercmd.PbObj) error{
	data, flag, err := common.EncodeGoCmd(uint16(cmd), msg)
	if err != nil {
		glog.Error("[Room] Send failed.", cmd, ", len: ",len(data), " err: ", err)
		return nil
	}

	this.AsyncSend(data, flag)
	return nil
}