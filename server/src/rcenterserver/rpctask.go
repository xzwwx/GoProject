package main

import (
	"common"
	"errors"
	"glog"
	"time"
)

var ErrFailed = errors.New("Failed.")

type RPCTask int


func (this *RPCTask) GetFreeRoom (args *common.ReqRoom, reply *common.RetRoom) error {

	room := FreeRoomMgr_GetMe().GetRoom(args.UserId)
	if room == nil {
		serverid, serveraddr, newsync := ServerTaskMgr_GetMe().GetServer()
		if serverid == 0 {
			return ErrFailed
		}
		reply.ServerId = serverid
		reply.Address = serveraddr
		reply.RoomId = RoomIdMgr_GetMe().GenerateId()
		reply.EndTime = uint32(time.Now().Unix() + int64(600 * time.Second))
		reply.NewSync = newsync
		FreeRoomMgr_GetMe().AddRoom(serverid, serveraddr, newsync, reply.RoomId, reply.EndTime)
	} else {
		reply.ServerId = room.ServerId
		reply.Address = room.Address
		reply.RoomId = room.RoomId
		reply.EndTime = room.EndTime
		reply.NewSync = room.NewSync

		// Load weight increment
		// To do

	}
	glog.Info("[RPC] Free Room ", args.UserId, ", ", reply.RoomId, ", ", reply.Address, ", ", *args, ",")
	return nil
}