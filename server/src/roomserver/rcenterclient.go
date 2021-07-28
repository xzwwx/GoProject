package main

import (
	"base/env"
	"common"
	"github.com/golang/protobuf/proto"
	"glog"
	"gonet"
	"usercmd"
)

type RCenterClient struct {
	task *gonet.TcpTask
	Id 	 uint16
}

var lclientm *RCenterClient

func RCenterClient_GetMe() *RCenterClient {
	if lclientm == nil {
		lclientm = &RCenterClient{
		}
	}
	return lclientm
}

func (this *RCenterClient) Connect() bool {

	loginaddr := env.Get("room","rcenter")
	mclient := &gonet.TcpClient{}
	conn, err := mclient.Connect(loginaddr)
	if err != nil {
		glog.Error("[Start] Connecting failed. ", loginaddr)
		return false
	}

	task := gonet.NewTcpTask(conn)
	task.Derived = this
	task.Start()

	this.task = task

	this.SendCmd(usercmd.CmdType_Login, &usercmd.ReqServerLogin{
		Address: (env.Get("room", "local")),
		Key: 	 (env.Get("room","key")),
		SerType: (common.ServerTypeRoom),
		NewSync: (true),
	})

	glog.Info("[Start] Connected to server successfully.", loginaddr)
	return true
}


func (this *RCenterClient) ParseMsg(data []byte, flag byte) bool {



	return true
}

func (this *RCenterClient) OnClose() {

}



/////////////////////////////////////
func (this *RCenterClient) SendCmd(cmd usercmd.CmdType, msg common.Message) bool {		//common.Message is an interface which is used ot Encode/Decode message
	//Encode
	data, flag, err, := common.EncodeGoCmd(uint16(cmd), msg)		///-----------------------
	if err != nil{
		glog.Info("[Service] Send Failed.", cmd, ", len:", len(data), ", err:", err)
		return false
	}

	this.task.AsyncSend(data, flag)

	return true
}
/////////////////////////////////////

func (this *RCenterClient) SendCmdToServer(serverid uint16, cmd usercmd.CmdType, msg proto.Message) bool {
	data, flag, err := common.EncodeCmd(uint16(cmd), msg)
	if err != nil{
		glog.Info("[Service] Send Failed.", cmd, ", len:", len(data), ", err:", err)
		return false
	}
	reqCmd := &usercmd.S2SCmd{
		ServerId: uint32(serverid),
		Flag: (uint32(flag)),
		Data: data,
	}
	return this.SendCmd(usercmd.CmdType_S2S, reqCmd)
}
/////////////////////////////////////////////


func (this *RCenterClient) GetId() uint16{
	return this.Id
}

func (this *RCenterClient) AddRoom(roomtype, roomid, endtime uint32, robot uint32) bool {
	reqCmd := &usercmd.ReqAddRoom{
		RoomType: roomtype,
		RoomId: roomid,
		EndTime: endtime,
	}
	return this.SendCmd(usercmd.CmdType_AddRoom, reqCmd)
}
func (this *RCenterClient) RemoveRoom(roomtype, roomid uint32, iscustom bool) bool {
	reqCmd := &usercmd.ReqRemoveRoom{
		//RoomType: roomtype,
		RoomId: roomid,
		IsCustom: iscustom,
	}
	return this.SendCmd(usercmd.CmdType_RemoveRoom, reqCmd)
}

//func (this *RCenterClient) UpdateRoom(roomtype, roomid, endtime uint32,  iscustom bool, robot uint32) bool {
//	reqCmd := &usercmd.Req{
//		RoomType: roomtype,
//		RoomId: roomid,
//		EndTime: endtime,
//	}
//	return this.SendCmd(usercmd.CmdType_AddRoom, reqCmd)
//}

func (this *RCenterClient) EndGame(roomid uint32, userid uint64) bool {
	reqCmd := &usercmd.ReqEndGame{
		RoomId: roomid,
		UserId: userid,
	}
	return this.SendCmd(usercmd.CmdType_EndGame, reqCmd)
}

func (this *RCenterClient) UpdateServer(roomnum, usernum int32) bool {
	reqCmd := &usercmd.ReqUpdateServer{
		RoomNum: (uint32(roomnum)),
		UserNum: (uint32(usernum)),

	}
	return this.SendCmd(usercmd.CmdType_UpdateServer, reqCmd)
}

