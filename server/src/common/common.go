package common

import (
	//"encoding/binary"
	//"glog"
	//"net/http"
	//"net/url"
	//"strconv"
	//"usercmd"
)

//Server Type
const(
	ServerTypeRoom	= 1 	// Room Server
	ServerTypeTeam 	= 2
	ServerTypeLogin	= 3		// Login Server
)

// Cmd size
const(
	CmdHeaderSize = 2
	ServerCmdSize = 1
	ServerIdSize = 4
	SubCmdSize = 2
)


// Player room token data
type UserData struct{
	ServerId	uint16 `redis:"ServerId"` 	// Vreify server Id
	Id 			uint64 `redis:"Id"`			// userid
	Account 	string `redis:"Account"` 	// username

	RoomId 		uint32 `redis:"RoomId"`
	RoomAddr 	string `redis:"RoomAddr"`

}

// Request message in game//////
type ReqMoveMsg struct {
	UserId uint64
	Speed uint32
	Direction uint32
}

type ReqLayBombMsg struct {
	UserId 	uint64
	X 		uint32
	Y 		uint32
}

// chi daoju
type ReqTriggerObjectMsg struct {
	UserId 	uint64
	ObjId	uint32		// daoju id
}

// bei zha
type ReqTriggerBombMsg struct {
	UserId 	uint64
	X		uint32
	Y 		uint32
}

// kill player
type ReqKillMsg struct{
	UserId uint64
	beKilled uint64
}
//////////////////

// Return from server


//Get Cmd
func GetCmd(buf []byte)uint16{
	if len(buf) <CmdHeaderSize{
		return 0
	}
	return uint16(buf[0])|uint16(buf[1])<<8
}