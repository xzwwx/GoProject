package main

import (
	"common"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type ScenePlayer struct {
	id 			uint64
	name 		string
	key 		string		// login key

	self 		*PlayerTask
	playertask 	*PlayerTask
	scene 		*Scene
	curflag 	map[uint32]*ScenePlayer
	otherPlayers 	map[uint64]*ScenePlayer
	rangebombs	[]*Bomb
	bombs		[]*Bomb

	senddie		bool

	isMove		bool
	X 			uint32
	Y 			uint32


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
		X: r.Uint32(),
		Y: r.Uint32(),
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



// Update players in scene
func (this *ScenePlayer) UpdateViewPlayers(scene *Scene){

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



func (this *ScenePlayer) sendSceneMsg(){

}