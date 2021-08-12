package main

import (
	"common"
	"fmt"
	"glog"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"usercmd"
)

type ScenePlayer struct {
	id   uint64
	name string
	key  string // login key

	self *PlayerTask
	// playertask *PlayerTask
	scene *Scene

	// curflag      map[uint32]*ScenePlayer
	otherPlayers map[uint64]*ScenePlayer
	rangeBombs   []*Bomb //current bombs
	// bombs        []*Bomb

	senddie bool

	// neng fou fen li chu lai
	PlayerMove
	isMove    int
	pos       Vector2
	nextpos   Vector2
	speed     uint32
	power     uint32
	direction uint32
	lifeState uint32 // 0: dead   1: alive 	2: jelly
	bombNum   uint32 //
	bombLeft  uint32

	//offline data
	udata *common.UserData

	movereq    *common.ReqMoveMsg
	laybombreq *common.ReqLayBombMsg
	objreq     *common.ReqTriggerObjectMsg // chi daoju
	killreq    *common.ReqKillMsg
}

func NewScenePlayer(player *PlayerTask, scene *Scene) *ScenePlayer {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := &ScenePlayer{
		id:           player.id,
		scene:        scene,
		self:         player,
		otherPlayers: make(map[uint64]*ScenePlayer),
		senddie:      false,
		pos: Vector2{
			x: r.Float64(),
			y: r.Float64(),
		},
	}
	return s
}

//Sync from scene
func (this *ScenePlayer) Update(perTime float64, scene *Scene) {
	//update move
	task := this.self
	if task != nil && atomic.CompareAndSwapInt32(&task.hasMove, 1, 0) {
		this.Move(scene, float64(atomic.LoadInt32(&task.speed)), float64(atomic.LoadInt32(&task.direction)))
	}

	//update bomb
	//this.rangebombs = this.rangebombs[:0]
	//this.rangebombs = append(this.rangebombs,this.bombs)
	for _, ball := range this.rangeBombs {
		if ball.isdelete == 1 {
			continue
		}
		scene.rangeBalls = append(scene.rangeBalls, ball)

		now := scene.now.Unix()
		if ball.layTime < now {
			continue
		}

	}

}

// Send update msg to all
func (this *ScenePlayer) sendSceneMsg(scene *Scene) {
	var (
		Moves       = scene.pool.moveFree[:0]
		ScenePlayer = scene.pool.msgPlayer[:0]
		MsgBomb     = scene.pool.msgBomb[:0]
		//Bombs =
	)
	// players move msg
	for _, playermove := range this.otherPlayers {
		if playermove.isMove == 1 {
			move := &usercmd.MsgPlayerMove{
				Id: playermove.id,
				X:  int32(playermove.pos.x),
				Y:  int32(playermove.pos.y),
				Nx: int32(playermove.nextpos.x),
				Ny: int32(playermove.nextpos.y),
			}
			Moves = append(Moves, move)
		}
		sp := &usercmd.ScenePlayer{
			Id:      playermove.id,
			BombNum: playermove.bombNum,
			Power:   playermove.power,
			Speed:   playermove.speed,
			State:   playermove.lifeState,
			X:       float32(playermove.pos.x),
			Y:       float32(playermove.pos.y),
			IsMove:  true,
		}
		ScenePlayer = append(ScenePlayer, sp)
	}

	for _, bomb := range this.rangeBombs {
		mb := &usercmd.MsgBomb{
			X:        bomb.pos.X,
			Y:        bomb.pos.Y,
			IsDelete: bomb.isdelete,
			Power:    bomb.power,
		}
		MsgBomb = append(MsgBomb, mb)
	}

	if len(Moves) != 0 {
		msg := &scene.pool.msgScene
		msg.Moves = Moves
		msg.Frame = scene.frame
		msg.Bombs = MsgBomb
		msg.Players = ScenePlayer
		msg.Id = this.id
		msg.X = float32(this.pos.x)
		msg.Y = float32(this.pos.y)

		this.sendSceneMsgToNet(msg, scene)
	}

}

// 同步场景信息
func (this *ScenePlayer) sendSceneMsgToNet(msg *usercmd.MsgScene, scene *Scene) {
	if this.self != nil {
		// 1.通过msgEncode编码发送场景信息
		newPos := msgSceneToBytes(uint16(usercmd.MsgTypeCmd_SceneSync), msg, scene.msgBytes)

		// 2.通过 Grpc 自带方法编码 发送场景信息
		data, ok := common.EncodeToBytes(uint16(usercmd.MsgTypeCmd_SceneSync), msg)
		if !ok {
			glog.Info("[玩家]更新场景失败 cmd:", uint16(usercmd.MsgTypeCmd_SceneSync))
			return
		}
		this.self.AsyncSend(data, 0)

		//TCP
		this.self.AsyncSend(scene.msgBytes[:newPos], 0)
	}
}

// Update players in scene
func (this *ScenePlayer) UpdateViewPlayers(scene *Scene) {
}

func (this *ScenePlayer) AsyncSend(buffer []byte, flag byte) {

	this.self.AsyncSend(buffer, flag)
}

// Time Action
func (this *ScenePlayer) TimeAction(room *Room, timenow time.Time) bool {

	return true
}

//Lay bomb
func (this *ScenePlayer) LayBomb(room *Room, x, y int32) {

	var (
		isLayBomb = false
		scene     = &room.Scene
	)
	if this.bombLeft > 0 {
		// cell := scene.GetCellState(uint32(this.pos.x), uint32(this.pos.y))
		cell := scene.GetCellState(uint32(x), uint32(y))
		if cell == 0 { // ke fang zha dan
			timenow := time.Now().Unix()
			bomb := &Bomb{
				pos: VectorInt{
					// X: int32(this.pos.x),
					// Y: int32(this.pos.y),
					X: x,
					Y: y,
				},
				player:   this,
				layTime:  timenow,
				isdelete: 0,
			}
			this.bombLeft--

			// this.rangeBombs = append(this.rangeBombs, bomb)
			// scene.rangeBalls = append(scene.rangeBalls, bomb)

			// 另一种形式管理炸弹
			room.bombmgr.Add(bomb)

			isLayBomb = true
		}

	}
	if isLayBomb {
		this.SendCmd(usercmd.MsgTypeCmd_LayBomb, &usercmd.MsgLayBomb{})
	}

	if room != nil {

	}

}

// Player Send Cmd
func (this *ScenePlayer) SendCmd(cmd usercmd.MsgTypeCmd, msg common.Message) bool {
	data, ok := common.EncodeToBytes(uint16(cmd), msg)
	if !ok {
		glog.Info("[Player] Send cmd:", cmd, ", len:", (len(data)))
		return false
	}
	this.AsyncSend(data, 0)
	return true
}

// 发送状态
func (this *ScenePlayer) SendState(room *Room) {
	room.BroadcastMsg(usercmd.MsgTypeCmd_PlayerState, PlayerStateCmd(this.id, int32(this.lifeState)))
}

//////////ScenePlayer manager//////////
type ScenePlayerMgr struct {
	mutex   sync.RWMutex
	players map[string]*ScenePlayer
}

var sptaskm *ScenePlayerMgr

func ScenePlayerMgr_GetMe() *ScenePlayerMgr {
	if sptaskm == nil {
		sptaskm = &ScenePlayerMgr{
			players: make(map[string]*ScenePlayer),
		}
	}
	return sptaskm
}

// Add
func (this *ScenePlayerMgr) Add(task *ScenePlayer) {
	this.mutex.RLock()
	this.players[task.key] = task
	this.mutex.RUnlock()
}

// Get player by key
func (this *ScenePlayerMgr) GetPlayer(key string) *ScenePlayer {
	this.mutex.RLock()
	player, _ := this.players[key]
	this.mutex.RUnlock()
	return player
}

// 删除场景玩家
func (this *ScenePlayerMgr) Removes(splayers map[uint64]*ScenePlayer) {
	this.mutex.Lock()
	for _, player := range splayers {
		delete(this.players, player.key)
	}
	fmt.Println("删除场景玩家")
	this.mutex.Unlock()
}

////////////
func (this *ScenePlayer) Move(scene *Scene, speed, direction float64) {

	this.isMove = 1
	this.CaculateNext(direction)          // 计算下一个位置
	this.scene.BorderCheck(&this.nextpos) // 保证计算得到的下一位置不超出地图范围

	// TODO 判断移动路径上是否有障碍物
	cellstate := scene.GetCellState(uint32(this.nextpos.x), uint32(this.nextpos.y))
	if cellstate == 0 {
		this.MoveVec(scene, speed, direction)
	} else if cellstate == 3 {
		// 吃道具
		this.MoveVec(scene, speed, direction)
		obj := scene.GetObjType(int32(this.nextpos.x), int32(this.nextpos.y))
		switch obj { // 1: 加速 speed   2：威力 power  3：数量  Bombnum
		case 1:
			this.speed++
			atomic.StoreInt32(&scene.gameMap.gamemap[int32(this.nextpos.x)][int32(this.nextpos.y)], 0)
		case 2:
			this.power++
			atomic.StoreInt32(&scene.gameMap.gamemap[int32(this.nextpos.x)][int32(this.nextpos.y)], 0)
		case 3:
			this.bombNum++
			atomic.StoreInt32(&scene.gameMap.gamemap[int32(this.nextpos.x)][int32(this.nextpos.y)], 0)
		}
	} else {
		// 位置不改变
	}
}

func (this *ScenePlayer) MoveVec(scene *Scene, speed, direction float64) {
	this.pos = this.nextpos
}

// 计算下一个位置
func (this *ScenePlayer) CaculateNext(direction float64) {
	this.nextpos.x = this.pos.x + float64(this.speed)*(math.Cos(direction*math.Pi/2))*0.04
	this.nextpos.y = this.pos.y + float64(this.speed)*(math.Sin(direction*math.Pi/2))*0.04

}
