package main

type Bomb struct {
	id
	player *ScenePlayer
	layTime 	int64
	pos 		Pos
	lastTime 	int64 	// duojiu zha
	isdelete	bool 	// zha le ma
	boxDistroiedList 	[]*Obstacle
}
