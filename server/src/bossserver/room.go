package main

import (
	"glog"
	"math"
	"sync"
	"time"
)

type Room struct {
	Scene
	//mutex    sync.RWMutex
	id          uint32                 //房间id
	roomType    uint32                 //房间类型
	//players     map[uint64]*PlayerTask //房间内的玩家
	curNum      uint32                 //当前房间内玩家数
	bombCount 	uint32
	isStart     bool
	timeLoop    uint64
	stopCh      chan bool
	isStop      bool
	iscustom	bool
	frame 		uint32

	// Operation
	scene 			*Scene
	opChan			chan *opMsg 	// player operation msg
	chan_PlayerOp 	chan *PlayerOp
	chan_Control 	chan int
	chan_AddPlayer 		chan *PlayerTask
	chan_RemovePlayer 	chan *PlayerTask


	newLogin 	map[uint64]bool
	loginHis	map[uint64]bool
	newLoginMutex sync.Mutex

	preFrameTime 	time.Time

	msgBytes 	[]byte

	now 		time.Time
	startTime 	time.Time
	max_num		uint32			//max player number. default :8
	totalTime 	uint64 //in second
	endchan     chan bool
}

type opMsg struct {
	op uint32
	args interface{}
}

//type PlayerOpt struct {
//	pTask 		*PlayerTask
//	pPlayer 	*ScenePlayer
//	Opts		*UserOpt
//}

func NewRoom(rtype, rid uint32, player *PlayerTask) *Room{
	room := &Room{

	}
	return room
}

func (this *Room) Start()bool{




	return true
}

func (this *Room) Stop() bool {



	return true
}

func (this *Room) IsClosed() bool {



	return true
}


func (this *Room) AddLoginUser(UID uint64) (result bool){
	this.newLoginMutex.Lock()
	defer this.newLoginMutex.Unlock()

	result = this.loginHis[UID]
	if !result {
		this.loginHis[UID] = true
	}
	this.newLogin[UID] = true
	return
}





// Main game loop
func (this *Room) Loop() {
	timeTicker := time.NewTicker(time.Millisecond * 20)
	stop := false
	defer func(){
		this.Stop()
		timeTicker.Stop()
		RoomMgr_GetMe().RemoveRoom(this)
	}()

	for {
		this.now = time.Now()
		select {
		case <-timeTicker.C:
			// 0.02s
			if this.timeLoop % 2 == 0 {
				this.Update(0.04)
			}
			//0.1s
			if this.timeLoop % 5 == 0{
				this.frame ++

				this.SendRoomMsg()
			}

			//1s
			if this.timeLoop % 100 == 0 {

				//this.scene.sendTime(this.totalTime - this.timeLoop/100)
			}
			if this.timeLoop != 0 && this.timeLoop % (this.totalTime * 100) == 0 {
				stop = true
			}
			this.timeLoop ++

			if this.isStop{
				stop = true
			}
			case op := <-this.chan_PlayerOp:
				//this.scene.UpdateOP(op)
				switch op.opType {
				case PlayerLayBombOp:
					this.LayBomb(op.playerId)
				case PlayerMoveOp:

				}
		}
	}
	this.Close()
}

func (this *Room) Close(){
	if !this.isStop {
		this.scene.SendOverMsg()
		this.isStop = true
		RoomMgr_GetMe().endChan <- this.id
	}
}

func(this *Room) Update(per float64) {
	//this.scene.UpdatePos()
	starttime := time.Now()
	ftime := starttime.Sub(this.preFrameTime).Milliseconds()
	this.preFrameTime = starttime

	this.UpdatePlayers(per)
	rtime := time.Now().Sub(starttime).Milliseconds()
	if math.Abs(float64(ftime - 40)) > 20 || rtime > 20 {
		glog.Info("[Statistic] State sync.", this.roomType, ", ", this.id, ", ", this.frame, ", ", ftime, ", ", rtime)
	}
}

func (this *Room) TimeAction() {

}

//Send cmd to room pthread
func (this *Room) Control(ctrl int) bool {
	if this.IsClosed() {
		return false
	}
	this.chan_Control <- ctrl
	return true

}

//Lay bomb
func (this *Room) LayBomb(playerId uint64) {
	player, ok := this.players[playerId]
	if !ok {
		return
	}
	player.LayBomb(this)

}

