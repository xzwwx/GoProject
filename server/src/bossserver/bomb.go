package main

type Bomb struct {
	id uint32
	player *ScenePlayer
	layTime 	int64
	pos 		Vector2
	lastTime 	int64 	// duojiu zha
	isdelete	bool 	// zha le ma
	boxDistroiedList 	[]*Obstacle
}

func (this *Bomb) TimeAction(){

}
