package main

type FreeRoom struct {
	room *Room
}

func NewFreeRoom(room *Room) *FreeRoom {
	return &FreeRoom{
		room: room,
	}
}

func (this *FreeRoom) LoadData() bool {

	return true
}

func (this *FreeRoom) Init() bool {

	return true
}

// Add player
//func (this * FreeRoom) AddPlayer(player *PlayerTask, p *ScenePlayer) (*ScenePlayer, /* *TeamMems */bool) {

//}

// Time Action
func (this *FreeRoom) TimeAction() {

}

//Sync player data
func (this *FreeRoom) SyncData() bool {
	return true
}

// End game
func (this *FreeRoom) DoEndGame() {

}

//Depose Player end game
func (this *FreeRoom) DoUserEndGame(player *ScenePlayer) {

}
