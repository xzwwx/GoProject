package main

import (
	"usercmd"
)

//player pool
type BallPool struct {

	moveFree 	[]*usercmd.BallMove
	msgPlayer 	[]*usercmd.MsgPlayer
	msgScene 	usercmd.MsgScene

}
