package main

import (
	"glog"
	"runtime/debug"
	"sync"
	"time"
)

type RoomMgr struct {
	runmutex sync.RWMutex
	runrooms map[uint32]*Room
	endmutex sync.RWMutex
	endrooms map[uint32]int64
	endChan  chan uint32
}

var roommgr *RoomMgr

func RoomMgr_GetMe() *RoomMgr {
	if roommgr == nil {
		roommgr = &RoomMgr{
			runrooms: make(map[uint32]*Room),
			endrooms: make(map[uint32]int64),
		}
		roommgr.Init()
	}
	return roommgr
}

func (this *RoomMgr) Init() {
	go func() {
		mintick := time.NewTicker(time.Minute)
		defer func() {
			if err := recover(); err != nil {
				glog.Error("[Exception] Error ", err, "\n", string(debug.Stack()))
			}
			mintick.Stop()
		}()

		for {
			select {
			case <-mintick.C:
				this.ChkEndRoomId()
			}
		}
	}()
}

// Add end room
func (this *RoomMgr) AddEndRoomId(roomid uint32) {
	this.endmutex.Lock()
	this.endrooms[roomid] = time.Now().Unix() + MAX_KEEPEND_TIME
	this.endmutex.Unlock()
}

// shifoushi End room
func (this *RoomMgr) IsEndRoom(roomid uint32) bool {
	this.endmutex.Lock()
	defer this.endmutex.Unlock()
	endtime, ok := this.endrooms[roomid]
	if !ok {
		return false
	}
	if endtime < time.Now().Unix() {
		delete(this.endrooms, roomid)
		return false
	}
	return true
}

// Check endroom list
func (this *RoomMgr) ChkEndRoomId() {
	timenow := time.Now().Unix()
	this.endmutex.Lock()
	for roomid, endtime := range this.endrooms {
		if endtime > timenow {
			continue
		}
		delete(this.endrooms, roomid)
	}
	this.endmutex.Unlock()
}

// Add room
func (this *RoomMgr) AddRoom(room *Room) (*Room, bool) {
	this.runmutex.Lock()
	defer this.runmutex.Unlock()
	oldroom, ok := this.runrooms[room.id]
	if ok {
		return oldroom, true
	}
	this.runrooms[room.id] = room
	return room, true
}

// Delete running room
func (this *RoomMgr) RemoveRoom(room *Room) {
	this.runmutex.Lock()
	delete(this.runrooms, room.id)
	this.runmutex.Unlock()
	this.AddEndRoomId(room.id)
	RCenterClient_GetMe().RemoveRoom(room.roomType, room.id, room.iscustom)
	RCenterClient_GetMe().UpdateServer(this.getNum(), PlayerTaskMgr_GetMe().GetNum())
	glog.Info("[Room] Remove Room[", room.roomType, ". ", room.id)
}

func (this *RoomMgr) getNum() (roomnum int32) {
	this.runmutex.Lock()
	roomnum = int32(len(this.runrooms))
	this.runmutex.Unlock()
	return
}

// Create room
func (this *RoomMgr) NewRoom(rtype, rid uint32, player *PlayerTask) *Room {
	room, ok := this.AddRoom(NewRoom(rtype, rid, player))
	if ok {
		if !room.Start() {
			this.RemoveRoom(room)
			return nil
		}

	}
	return room
}

// Get rooms
func (this *RoomMgr) GetRooms() (rooms []*Room) {
	this.runmutex.RLock()
	for _, room := range this.runrooms {
		rooms = append(rooms, room)
	}
	this.runmutex.RUnlock()
	return
}

//Get room by id
func (this *RoomMgr) getRoomById(rid uint32) *Room {
	this.runmutex.RLock()
	room, ok := this.runrooms[rid]
	if !ok {
		return nil
	}
	return room
}

func (this *RoomMgr) AddPlayer(player *PlayerTask) bool {

	room := this.getRoomById(player.udata.RoomId)

	if this.IsEndRoom(player.udata.RoomId) {
		glog.Error("[房间]已结束", player.udata.Id, ", ", player.udata.Account, ", ", player.udata.RoomId)
		return false
	}

	if room == nil {
		room = this.NewRoom(0, player.udata.RoomId, player)
		if room == nil {
			return false
		}
	}
	glog.Info("[房间] 自由模式 ", room.id, ", ", player.udata.Id, ", ", player.udata.Account)

	room.IncPlayerNum()

	player.room = room
	player.room.chan_AddPlayer <- player

	return true
}
