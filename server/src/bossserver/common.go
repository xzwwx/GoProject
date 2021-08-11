package main

import "usercmd"

type UserOpt struct {
	Id         uint8
	hasMove    bool
	hasBreak   bool
	hasCombine bool
	Power      uint8
	Speed      uint16
	PosX       uint16
	PosY       uint16
	OptSeq     uint8
}

// Obstacle
type Obstacle struct {
	Id    uint32
	pos   VectorInt
	OType uint32 //0:wall   1: ke zha
}

type VectorInt struct {
	X int32
	Y int32
}

type Vector2 struct {
	x float64
	y float64
}

func PlayerStateCmd(uid uint64, state int32) *usercmd.RetPlayerState {
	return &usercmd.RetPlayerState{
		UserId:    uid,
		UserState: state,
	}
}
