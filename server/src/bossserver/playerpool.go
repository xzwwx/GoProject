package main

import (
	"usercmd"
)

//player pool
type BallPool struct {
	moveFree  []*usercmd.MsgPlayerMove
	msgPlayer []*usercmd.MsgPlayer
	msgScene  usercmd.MsgScene
}
