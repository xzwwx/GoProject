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

	gameMap *Map
	objects *Objects
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
	rangeBalls []*Bomb //map
	// rangePlayers   []*ScenePlayer
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
	// this.rangePlayers = this.rangePlayers[:0]

	this.startTime = time.Now()

	this.Obstacle = GenerateRandMap()

	// 加载地图
	this.gameMap = &Map{} // 50x50

	// 初始化道具位置

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

// Check mapcell   0: null  1:wall  2:Obstacle 3: 道具
func (this *Scene) GetCellState(x, y uint32) int32 {
	return this.gameMap.gamemap[x][y]
}

// Check mapcell   0: null  1:wall  2:Obstacle 3: 道具
func (this *Scene) GetObjType(x, y int32) int32 {
	return this.objects.obj[x][y]
}

// 保证位置不超出地图范围
func (this *Scene) BorderCheck(pos *Vector2) {
	if pos.x < 0 {
		pos.x = 0
	} else if pos.x >= this.sceneWidth {
		pos.x = this.sceneWidth - 0.01
	}
	if pos.y < 0 {
		pos.y = 0
	} else if pos.y >= this.sceneHeight {
		pos.y = this.sceneHeight - 0.01
	}
}
