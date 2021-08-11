package main

import (
	"common"
	"glog"
	"math"
	"sync"
	"sync/atomic"
	"time"
	"usercmd"
)

type Room struct {
	Scene
	//mutex    sync.RWMutex
	id       uint32 //房间id
	roomType uint32 //房间类型
	//players     map[uint64]*PlayerTask //房间内的玩家
	playerNum int32 //当前房间内玩家数
	bombCount uint32
	isStart   bool
	timeLoop  uint64
	stopCh    chan bool
	isStop    bool
	iscustom  bool
	frame     uint32
	isclosed  int32

	rmode RoomMode

	// 处理炸弹
	bombmgr *BombMgr

	// Operation
	scene             *Scene
	opChan            chan *opMsg // player operation msg
	chan_PlayerOp     chan *PlayerOp
	chan_Control      chan int
	chan_AddPlayer    chan *PlayerTask
	chan_RemovePlayer chan *PlayerTask

	newLogin      map[uint64]bool
	loginHis      map[uint64]bool
	newLoginMutex sync.Mutex

	preFrameTime time.Time

	msgBytes []byte

	now       time.Time
	startTime time.Time
	max_num   uint32 //max player number. default :8
	totalTime uint64 //in second
	endchan   chan bool
}

type opMsg struct {
	op   uint32
	args interface{}
}

type RoomMode interface {
	LoadData() bool
}

//type PlayerOpt struct {
//	pTask 		*PlayerTask
//	pPlayer 	*ScenePlayer
//	Opts		*UserOpt
//}

func NewRoom(rtype, rid uint32, player *PlayerTask) *Room {
	room := &Room{}
	return room
}

func (this *Room) Start() bool {
	if !atomic.CompareAndSwapInt32(&this.isclosed, -1, 0) {
		return false
	}
	this.rmode = NewFreeRoom(this)

	this.bombmgr = NewBombMgr(this)
	this.Scene.Init(this)

	go this.Loop()
	glog.Info("[房间] 创建房间 ", ", ", this.id, ", ")
	return true
}

// 房间停止
func (this *Room) Stop() bool {
	if !atomic.CompareAndSwapInt32(&this.isclosed, 0, 1) {
		return false
	}
	this.destory()
	glog.Info("[房间] 销毁房间 ", this.id, ", ", len(this.players))

	return true
}

func (this *Room) IsClosed() bool {
	return atomic.LoadInt32(&this.isclosed) != 0
}

// 删除房间
func (this *Room) destory() {
	this.Stop()
	go func(room *Room) {
		ScenePlayerMgr_GetMe().Removes(room.players)

		// redis 清理玩家
		// TODO

	}(this)
	glog.Info("[房间] 结算完成", this.id, ", ", this.GetPlayerNum())
}

func (this *Room) AddLoginUser(UID uint64) (result bool) {
	this.newLoginMutex.Lock()
	defer this.newLoginMutex.Unlock()

	result = this.loginHis[UID]
	if !result {
		this.loginHis[UID] = true
	}
	this.newLogin[UID] = true
	return
}

// 玩家 +1
func (this *Room) IncPlayerNum() {
	atomic.AddInt32(&this.playerNum, 1)
}

// 返回玩家数
func (this *Room) GetPlayerNum() int32 {
	return atomic.LoadInt32(&this.playerNum)
}

// Main game loop
func (this *Room) Loop() {
	timeTicker := time.NewTicker(time.Millisecond * 20)
	//stop := false
	defer func() {
		this.Stop()
		timeTicker.Stop()
		RoomMgr_GetMe().RemoveRoom(this)
	}()

	for {
		this.now = time.Now()
		select {
		case <-timeTicker.C:
			// 0.02s
			if this.timeLoop%2 == 0 {
				this.Update(0.04)
			}
			//0.1s
			if this.timeLoop%5 == 0 {
				this.frame++

				this.SendRoomMsg()
				this.bombmgr.ExecAction()
			}

			//1s
			if this.timeLoop%100 == 0 {

				//this.scene.sendTime(this.totalTime - this.timeLoop/100)
			}
			if this.timeLoop != 0 && this.timeLoop%(this.totalTime*100) == 0 {
				// stop = true
			}
			this.timeLoop++

			if this.isStop {
				// stop = true
			}
		case op := <-this.chan_PlayerOp:
			//this.scene.UpdateOP(op)
			switch op.opType {
			case PlayerLayBombOp:
				this.LayBomb(op.playerId, op.x, op.y)
			case PlayerMoveOp:

			}
		}

	}
	//stop = true
	this.Close()
}

func (this *Room) Close() {
	if !this.isStop {
		//this.scene.SendOverMsg()
		this.isStop = true
		RoomMgr_GetMe().endChan <- this.id
	}
}

func (this *Room) Update(per float64) {
	//this.scene.UpdatePos()
	starttime := time.Now()
	ftime := starttime.Sub(this.preFrameTime).Milliseconds()
	this.preFrameTime = starttime

	this.UpdatePlayers(per)
	rtime := time.Now().Sub(starttime).Milliseconds()
	if math.Abs(float64(ftime-40)) > 20 || rtime > 20 {
		glog.Info("[Statistic] State sync.", this.roomType, ", ", this.id, ", ", this.frame, ", ", ftime, ", ", rtime)
	}
}

//广播消息
func (this *Room) BroadcastMsg(msgNo usercmd.MsgTypeCmd, msg common.Message) {
	data, ok := common.EncodeToBytes(uint16(msgNo), msg)
	if !ok {
		glog.Error("[广播] 发送消息失败 ", msgNo)
		return
	}
	for _, player := range this.Scene.players {
		player.AsyncSend(data, 0)
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
func (this *Room) LayBomb(playerId uint64, x, y int32) {
	player, ok := this.players[playerId]
	if !ok {
		return
	}
	player.LayBomb(this, x, y)

}
