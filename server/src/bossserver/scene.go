package main

import (
	"time"
)

const (
	RoomWidth  = 200
	RoomHeight = 200
	CellWidth  = 5
	CellHeight = 5

	RoomSize = 200
)

type Scene struct {
	players  map[uint64]*ScenePlayer
	room     *Room
	Obstacle *map[uint32]*Obstacle

	gameMap Map
	// Scene info
	sceneWidth  float64
	sceneHeight float64
	now         time.Time
	startTime   time.Time
	frame       uint32

	ballID uint32
	pool   *BallPool // Player Pool

	msgBytes []byte

	// temple things
	rangeBalls     []*Bomb //map
	rangePlayers   []*ScenePlayer
	rangeObstacles []*Obstacle
}

func NewScene(room *Room) *Scene {
	scene := &Scene{
		room:    room,
		players: make(map[uint64]*ScenePlayer),
	}
	scene.Init(room)
	return scene
}

func (this *Scene) Init(room *Room) {

	// 房间指针
	this.room = room

	this.players = make(map[uint64]*ScenePlayer)

	this.rangeBalls = this.rangeBalls[:0]
	this.rangePlayers = this.rangePlayers[:0]

	this.startTime = time.Now()

	this.Obstacle = GenerateRandMap()
	// Init map
	this.gameMap = Map{}
}

func (this *Scene) AddPlayer(p *PlayerTask) {
	this.players[p.id] = NewScenePlayer(p, this)
}

func (this *Scene) SendRoomMsg() {
	for _, p := range this.players {
		p.sendSceneMsg(this)
	}
}

func (this *Scene) UpdatePlayers(per float64) {

	// Depose player logic
	//if this.room.roomType
	//for i := 0; i < len(this.players); i++ {
	//	player, _ := this.players[]
	//}
	for _, player := range this.players {
		player.Update(per, this)
	}

}

// Check mapcell   0: null  1:wall  2:Obstacle
func (this *Scene) GetCellState(x, y uint32) int32 {

	return this.gameMap.gamemap[x][y]
}
