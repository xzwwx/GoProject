package main

import (
	"errors"
	"fmt"
	"glog"
	"time"
	"usercmd"
)

var ErrFailed = errors.New("Failed.")

type RPCTask int

func (this *RPCTask) GetFreeRoom(args *usercmd.ReqIntoRoom, reply *usercmd.RetIntoRoom) error {
	fmt.Println("Get room [start]", *args.UId)
	room := FreeRoomMgr_GetMe().GetRoom(*args.UId)
	if room == nil {
		fmt.Println("No room")
		serverid, serveraddr, newsync := ServerTaskMgr_GetMe().GetServer()
		fmt.Println("ServerId ", serverid, ", adderss: ", serveraddr)
		if serverid == 0 {
			return ErrFailed
		}
		// reply.ServerId = &serverid
		reply.Addr = &serveraddr
		reply.RoomId = RoomIdMgr_GetMe().GenerateId()
		fmt.Println("RoomId", reply.RoomId)
		// reply.EndTime = uint32(time.Now().Unix() + int64(600*time.Second))
		reply.NewSync = &newsync
		FreeRoomMgr_GetMe().AddRoom(serverid, serveraddr, newsync, *reply.RoomId, uint32(time.Now().Unix()+int64(600*time.Second)))
	} else {

		fmt.Println("Set reply...")

		// reply.ServerId = room.ServerId
		reply.Addr = &room.Address
		reply.RoomId = &room.RoomId
		// reply.EndTime = room.EndTime
		reply.NewSync = &room.NewSync

		// fmt.Println("ServerId: ", reply.ServerId, ", Addr:", reply.Address, ", RoomId", reply.RoomId, ", ", reply.EndTime, ", ", reply.NewSync)

		// Load weight increment
		// To do

	}
	glog.Info("[RPC] Free Room ", *args.UId, ", ", reply.RoomId, ", ", reply.Addr, ", ", *args, ",")
	return nil
}
