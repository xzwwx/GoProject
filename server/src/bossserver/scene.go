package main

import "time"

const (
	RoomWidth		= 200
	RoomHeight		= 200
	CellWidth		= 5
	CellHeight		= 5

	RoomSize		= 200
)

type Scene struct {
	players map[uint64]*ScenePlayer
	room 	*Room
	Obstacle *map[uint32]*Obstacle

	sceneWidth 		float64
	sceneHeight 	float64
	now 			time.Time
	startTime 		time.Time
	frame			uint32
	ballID			uint32

	// temple things
	rangeBalls 		[]*Bomb //map
	rangePlayers	[]*ScenePlayer
	rangeObstacles 	[]*Obstacle
}

func NewScene(room *Room) *Scene {
	scene := &Scene{
		room:room,
		players: make(map[uint64]*ScenePlayer),
	}
	scene.Init()
	return scene
}

func (this *Scene) Init(){
	this.Obstacle = GenerateRandMap()
}

func (this *Scene) AddPlayer (p *PlayerTask) {
	this.players[p.id] = NewScenePlayer(p,this)
}

func (this *Scene) SendRoomMsg() {
	for _, p := range this.players {
		p.sendSceneMsg()
	}
}

func (this *Scene) UpdatePlayers(per float64)  {

	// Depose player logic
	//if this.room.roomType
	//for i := 0; i < len(this.players); i++ {
	//	player, _ := this.players[]
	//}
	for _, player := range this.players {
		player.Update(per, this)
	}

}

func (this *Scene) SendRoomMsg() {

}