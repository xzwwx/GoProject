package main

import (
	"common"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"usercmd"
)

type ScenePlayer struct {
	id 			uint64
	name 		string
	key 		string		// login key

	self 		*PlayerTask
	playertask 	*PlayerTask
	scene 		*Scene

	curflag 		map[uint32]*ScenePlayer
	otherPlayers 	map[uint64]*ScenePlayer
	rangeBombs	[]*Bomb		//current bombs
	bombs		[]*Bomb

	senddie		bool

	// neng fou fen li chu lai
	PlayerMove
	isMove		bool
	pos 		Vector2
	nextpos 	Vector2
	speed 		uint32
	direction 	uint32
	lifeState  	uint32  // 0: dead   1: alive 	2: jelly
	bombNum 	uint32  //
	bombLeft	uint32

	//offline data
	udata 		*common.UserData

	movereq 	*common.ReqMoveMsg
	laybombreq  *common.ReqLayBombMsg
	objreq 		*common.ReqTriggerObjectMsg		// chi daoju
	killreq 	*common.ReqKillMsg

}

func NewScenePlayer(player *PlayerTask, scene *Scene) *ScenePlayer {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := &ScenePlayer{
		id: player.id,
		scene: scene,
		playertask: player,
		curflag: make(map[uint32]*ScenePlayer),
		senddie: false,
		pos: Vector2{
			x: r.Uint32(),
			y: r.Uint32(),
		},
	}
	return s
}

//Sync from scene
func (this *ScenePlayer) Update(perTime float64, scene *Scene){
	//update move
	task := this.self
	if task != nil && atomic.CompareAndSwapInt32(&task.hasMove, 1, 0){
		this.Move(scene, float64(atomic.LoadInt32(&task.speed)), float64(atomic.LoadInt32(&task.direction)))
	}

	//update bomb
	//this.rangebombs = this.rangebombs[:0]
	//this.rangebombs = append(this.rangebombs,this.bombs)
	for _, ball := range this.bombs {
		if ball.isdelete {
			continue
		}
		scene.rangeBalls = append(scene.rangeBalls, ball)

		now := scene.now.Unix()
		if ball.layTime < now{
			continue
		}

	}

}

// Send update msg to all
func (this *ScenePlayer) sendSceneMsg(scene *Scene){
var (
	Moves = scene.pool.moveFree[:0]
	//Bombs =
)
	// players move msg
	for _, playermove := range this.otherPlayers {
		if playermove.isMove {
			move := &usercmd.BallMove{
				Id: playermove.id,
				X: int32(playermove.pos.x),
				Y: int32(playermove.pos.y),
				Nx: int32(playermove.nextpos.x),
				Ny: int32(playermove.nextpos.y),
			}
			Moves = append(Moves,move)
		}
	}

	if len(Moves) != 0 {
		msg := &scene.pool.msgScene
		msg.Moves = Moves
		msg.Frame = scene.frame

		this.sendSceneMsgToNet(msg, scene)
	}

}

func (this *ScenePlayer) sendSceneMsgToNet(msg *usercmd.MsgScene, scene *Scene) {
	if this.self != nil {
		newPos := msgSceneToBytes(uint16(usercmd.MsgTypeCmd_NewScene), msg, scene.msgBytes)

		//TCP
		this.self.AsyncSend(scene.msgBytes[:newPos], 0)
	}
}

// Update players in scene
func (this *ScenePlayer) UpdateViewPlayers(scene *Scene){
}

func (this *ScenePlayer) AsyncSend(buffer []byte, flag byte)  {

	this.self.AsyncSend(buffer, flag)
}

// Time Action
func (this *ScenePlayer) TimeAction(room *Room, timenow time.Time) bool {


	return true
}


//Lay bomb
func (this *ScenePlayer) LayBomb (room *Room) {

	var (
		isLayBomb = false
		scene = &room.Scene
	)
	if this.bombLeft > 0 {
		cell := scene.GetCellState(this.pos.x, this.pos.y)
		if cell == 0 {
			timenow := time.Now().Unix()
			bomb := &Bomb{
				pos: this.pos,
				player: this,
				layTime: timenow,
				isdelete: false,
			}
			this.rangeBombs = append(this.rangeBombs, bomb)
			this.bombLeft --
			scene.rangeBalls = append(scene.rangeBalls, bomb)
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
func (this * ScenePlayer) SendCmd(cmd usercmd.MsgTypeCmd, msg common.Message) bool {

}






//////////ScenePlayer manager//////////
type ScenePlayerMgr struct {
	mutex sync.RWMutex
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
func(this *ScenePlayerMgr) Add(task *ScenePlayer) {
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







////////////
func (this *ScenePlayer) Move(scene *Scene, speed, direction float64) {

	// move how much
	//vec := &Vector2{
	//
	//}
	this.MoveVec(scene, speed, direction)
}
func (this *ScenePlayer) MoveVec(scene *Scene, speed, direction float64 ){
	//this.X =
}
