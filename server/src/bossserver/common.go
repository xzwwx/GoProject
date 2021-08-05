package main

type UserOpt struct {
	Id 		uint8
	hasMove	bool
	hasBreak	bool
	hasCombine	bool
	Power 		uint8
	Speed 		uint16
	PosX		uint16
	PosY		uint16
	OptSeq		uint8
}

// Obstacle
type Obstacle struct {
	Id 	uint32
	pos Pos
	OType 	uint32  //0:wall   1: ke zha
}

type Pos struct {
	X 	uint32
	Y	uint32
}

type Vector2 struct {
	x uint32
	y uint32
}

