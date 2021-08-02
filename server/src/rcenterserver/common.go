package main

import "sync/atomic"

const (
	MAX_ROOM_SNUM	= 110	// normal room max number

)
type Room struct {
	RServerId 	uint16	//Room server ID
	RAddress 	string 	//Room server address
	EndTime		uint32
	UserNum 	int32
	RobotNum	uint32
	NewSync		bool
}

type SortRoom struct {
	ServerId	uint16
	Address 	string
	RoomId		uint32
	EndTime		uint32
	UserNum 	int32
	RobotNum	uint32
	NewSync 	bool
}

//////Room unique index
type RoomIdMgr struct {
	serverid  	uint8
	uniqueid	uint32

}

var roomidm *RoomIdMgr

func RoomIdMgr_GetMe() *RoomIdMgr {
	if roomidm == nil {
		roomidm = &RoomIdMgr{
		}
	}
	return roomidm
}

func (this *RoomIdMgr) GenerateId() uint32 {
	return uint32(this.serverid) << 24 | atomic.AddUint32(&this.uniqueid, 1) % 0xffffff
}