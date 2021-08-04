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

// Player room token data
type UserData struct{
	ServerId	uint16 `redis:"ServerId"` 	// Vreify server Id
	Id 			uint64 `redis:"Id"`			// userid
	Account 	string `redis:"Account"` 	// username

	RoomId 		uint32 `redis:"RoomId"`
	RoomAddr 	string `redis:"RoomAddr"`

}