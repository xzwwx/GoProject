package main

import (
	"sync"
	"time"
)

type FRoom struct {
	Room

	RUserNum	int32
	MaxNum		int32
	RoomId		uint32
	StartTime 	uint32
	Members		map[uint64]bool
	FUserNum	uint32
}



// Add member
func (this *FRoom) Add(uid uint64) bool {
	if this.Members == nil {
		this.Members = make(map[uint64]bool)
	}
	_, ok := this.Members[uid]
	if !ok {
		this.Members[uid] = true
		this.FUserNum += 1
	}
	return true
}

//Get Numbers
func (this *FRoom) MemNum() uint32 {
	if this.FUserNum < uint32(this.RUserNum) {
		return uint32(this.RUserNum)
	}
	return this.FUserNum
}



type FreeRoomMgr struct {
	mutex 	sync.RWMutex
	rooms 	map[uint32]*FRoom
	waitmutex 	sync.RWMutex
	waitrooms 	map[uint32]*FRoom
}

var froomm *FreeRoomMgr

func FreeRoomMgr_GetMe() *FreeRoomMgr {
	if froomm == nil {
		froomm = &FreeRoomMgr{
			rooms : make(map[uint32]*FRoom),
			waitrooms : make(map[uint32]*FRoom),
		}
	}
	return froomm
}

//Add Room
func (this *FreeRoomMgr) AddRoom(rserverid uint16, raddress string, newsync bool, roomid uint32, endtime uint32) *FRoom {
	this.mutex.Lock()
	room := &FRoom{
		Room:Room{
			RServerId: rserverid,
			RAddress: raddress,
			EndTime: endtime,
			UserNum: 1,
			NewSync: newsync,
		},
		RUserNum: 1,
		RoomId: roomid,
		StartTime: uint32(time.Now().Unix()),
	}
	this.rooms[roomid] = room
	this.mutex.Unlock()
	return room
}

func (this *FreeRoomMgr) RemoveRoom(roomid uint32){
	this.mutex.Lock()
	delete(this.rooms, roomid)
	this.mutex.Unlock()
}

func (this *FreeRoomMgr) UpdateRoom(roomid uint32, usernum, rusernum int32) {
	this.mutex.Lock()
	room, ok := this.rooms[roomid]
	if ok {
		room.UserNum = usernum
		if usernum > rusernum{
			rusernum = usernum
		}
		if rusernum > room.RUserNum {
			room.RUserNum = rusernum
		}
	}

	this.mutex.Unlock()
}

func (this *FreeRoomMgr) IncRoomNum(userid uint64, roomid uint32) (usernum int32) {
	this.mutex.Lock()
	room, ok := this.rooms[roomid]
	if ok {
		usernum = room.RUserNum
		room.UserNum ++
		room.RUserNum ++
	}
	this.mutex.Unlock()
	return
}

func (this *FreeRoomMgr) GetRoom(userid uint64) *SortRoom {
	froom := &SortRoom{}
	nowTime := uint32(time.Now().Unix())
	//totalnum := ServerTaskMgr_GetMe().GetTotalNum()
	this.mutex.Lock()
	//roomnum := len(this.rooms)
	for rid, room := range this.rooms {
		if room.EndTime <= nowTime {
			delete(this.rooms, rid)
		}
		if room.RUserNum >= MAX_ROOM_SNUM {		// Full
			continue
		}

		if froom.RoomId == 0 || room.UserNum < froom.UserNum {
			froom.ServerId = room.RServerId
			froom.Address = room.RAddress
			froom.RoomId = room.RoomId
			froom.EndTime = room.EndTime
			froom.UserNum = room.UserNum
			froom.NewSync = room.NewSync
		}
	}
	this.mutex.Unlock()
	if froom.RoomId == 0 {
		return nil
	}
	this.IncRoomNum(userid, froom.RoomId)
	return froom
}

